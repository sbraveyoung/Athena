package easybits

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/SmartBrave/Athena/easyerrors"
	"github.com/SmartBrave/Athena/easyio"
)

type IntegerType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func Marshal(v interface{}, writer easyio.EasyWriter) (err error) {
	//TODO
	return nil
}

func Unmarshal(reader easyio.EasyReader, v interface{}) (err error) {
	tYpe := reflect.TypeOf(v)
	if tYpe.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid type of v:%+v", tYpe.Kind())
	}
	tYpe = tYpe.Elem()
	if tYpe.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type of *v:%+v", tYpe.Kind())
	}
	value := reflect.ValueOf(v).Elem()

	var b []byte
	for i := 0; i < tYpe.NumField(); i++ {
		fieldType := tYpe.Field(i)
		if !fieldType.IsExported() {
			continue
		}

		fieldKind := fieldType.Type.Kind()
		if fieldKind == reflect.Array || fieldKind == reflect.Slice {
			//array or slice
			if fieldElemKind := fieldType.Type.Elem().Kind(); fieldElemKind < reflect.Int || fieldElemKind > reflect.Uint64 {
				// continue
				panic("support integer type only!")
			}
		} else if fieldKind < reflect.Int || fieldKind > reflect.Uint64 {
			// continue
			panic("support integer type only!")
		}

		bitTag := fieldType.Tag.Get("bits")
		if bitTag == "" {
			continue
		}

		startByte, startBit, endByte, endBit, itemBit, err := parse(bitTag)
		if err != nil {
			return fmt.Errorf("invlid tag of bits:%s, err:%v", bitTag, err)
		}

		if fieldKind != reflect.Slice && fieldKind != reflect.Array && itemBit != 0 {
			return fmt.Errorf("invalid tag of bits:%s, not slice.", fieldType.Name)
		}

		bits := (endByte-startByte)*8 + (endBit - startBit)
		if bits == 0 {
			continue
		}

		if fieldKind == reflect.Slice {
			//nothing
		} else if fieldKind == reflect.Array {
			len := fieldType.Type.Len()
			if bits > itemBit*len {
				return fmt.Errorf("The bits of field %s is %d, but you want to store %d bit.", fieldType.Name, itemBit*len, bits)
			}
		} else {
			fieldBits := fieldType.Type.Bits()
			if bits > fieldBits {
				// bits = fieldBits
				return fmt.Errorf("The field %s's type %s is too small to store the num by tag:%d", fieldType.Name, fieldType.Type.Name(), bits)
			}
		}

		toReadBytes := 0
		if endBit == 0 && endByte > len(b) {
			toReadBytes = endByte - len(b)
		} else if endBit > 0 && endByte >= len(b) {
			toReadBytes = endByte - len(b) + 1
		}
		if toReadBytes > 0 {
			tmpBuf, err := reader.ReadN(uint32(toReadBytes))
			if err != nil {
				return fmt.Errorf("not enough data from reader")
			}
			b = append(b, tmpBuf...)
		}

		fieldValue := value.Field(i)
		if fieldKind != reflect.Slice && fieldKind != reflect.Array {
			var n uint64
			for byteIndex := startByte; byteIndex <= endByte; byteIndex++ {
				lbit, rbit := 0, 8
				if byteIndex == startByte {
					lbit = startBit
				}
				if byteIndex == endByte {
					rbit = endBit
				}
				for bitIndex := lbit; bitIndex < rbit; bitIndex++ {
					n <<= 1
					n |= uint64((b[byteIndex] >> (8 - bitIndex - 1)) & 0x01)
				}
			}

			if fieldValue.CanSet() {
				if fieldValue.CanInt() {
					fieldValue.SetInt(int64(n))
				} else if fieldValue.CanUint() {
					fieldValue.SetUint(n)
				}
			}
		} else {
			items := bits / itemBit
			if bits%itemBit != 0 {
				items++
			}
			var n uint64
			bit := 0
			index := 0
			for byteIndex := startByte; byteIndex <= endByte; byteIndex++ {
				lbit, rbit := 0, 8
				if byteIndex == startByte {
					lbit = startBit
				}
				if byteIndex == endByte {
					rbit = endBit
				}
				for bitIndex := lbit; bitIndex < rbit; bitIndex++ {
					n <<= 1
					n |= uint64((b[byteIndex] >> (8 - bitIndex - 1)) & 0x01)
					bit++
					if bit == itemBit || (byteIndex == endByte && bitIndex == rbit-1) {
						if fieldValue.CanSet() {
							if fieldKind == reflect.Slice { //TODO: optimization
								//growth
								if fieldValue.Len() <= index {
									fieldValue.Set(reflect.Append(fieldValue, reflect.New(fieldType.Type.Elem()).Elem()))
								}

								if elem := fieldValue.Index(index); elem.CanInt() {
									elem.SetInt(int64(n))
								} else if elem.CanUint() {
									elem.SetUint(n)
								}
							} else {
								if elem := fieldValue.Index(index); elem.CanInt() { //if index is out of range, panic
									elem.SetInt(int64(n))
								} else if elem.CanUint() {
									elem.SetUint(n)
								}
							}
						}
						bit = 0
						n = 0
						index++
					}
				}
			}
		}
	}
	return nil
}

var (
	invalidErr = fmt.Errorf("empty expression")
)

//[startByte.startBit,endByte.endBit)
func parse(expression string) (startByte, startBit, endByte, endBit, itemBit int, err error) {
	startByte, startBit, endByte, endBit, itemBit = 0, 0, 0, 0, 0
	if expression == "-" || expression == "" {
		return
	}

	if !strings.HasPrefix(expression, "[") || !strings.HasSuffix(expression, "]") {
		err = invalidErr
		return
	}

	var err0, err1, err2, err3, err4, err5 error
	slice := strings.Split(expression[1:len(expression)-1], ":")
	switch len(slice) {
	case 3:
		if slice[2] != "" {
			itemBit, err5 = strconv.Atoi(slice[2])
		} else {
			err0 = invalidErr
		}
		fallthrough
	case 2:
		if slice[1] != "" {
			endSlice := strings.Split(slice[1], ".")
			switch len(endSlice) {
			case 2:
				endByte, err3 = strconv.Atoi(endSlice[0])
				endBit, err4 = strconv.Atoi(endSlice[1])
			case 1:
				endByte, err3 = strconv.Atoi(endSlice[0])
				endBit = 0
			default:
				err0 = invalidErr
			}
		}

		if slice[0] != "" {
			startSlice := strings.Split(slice[0], ".")
			switch len(startSlice) {
			case 2:
				startByte, err1 = strconv.Atoi(startSlice[0])
				startBit, err2 = strconv.Atoi(startSlice[1])
			case 1:
				startByte, err1 = strconv.Atoi(startSlice[0])
				startBit = 0
			default:
				err0 = invalidErr
			}
		}
	case 1:
		if slice[0] != "" {
			startSlice := strings.Split(slice[0], ".")
			switch len(startSlice) {
			case 2:
				startByte, err1 = strconv.Atoi(startSlice[0])
				startBit, err2 = strconv.Atoi(startSlice[1])
				endByte, endBit = startByte, startBit+1
			case 1:
				//1
				startByte, err1 = strconv.Atoi(startSlice[0])
				startBit = 0
				endByte, endBit = startByte, 7
			default:
				err0 = invalidErr
			}
		}
	default:
		err = invalidErr
	}

	err = easyerrors.HandleMultiError(easyerrors.Simple(), err0, err1, err2, err3, err4, err5)
	if startByte < -1 || endByte < -1 || startBit < 0 || startBit > 7 || endBit < 0 || endBit > 7 || endByte < startByte || (endByte == startByte && endBit < startBit) {
		err = invalidErr
	}
	return
}
