package easybits

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/SmartBrave/Athena/easyerrors"
	"github.com/SmartBrave/Athena/easyio"
)

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
		field := tYpe.Field(i)
		typeFieldBits := field.Type.Bits()

		bitTag := field.Tag.Get("bits")
		if bitTag == "" {
			continue
		}

		startByte, startBit, endByte, endBit, err := parse(bitTag)
		if err != nil {
			return fmt.Errorf("invlid tag of bits:%s, err:%v", bitTag, err)
		}

		bits := (endByte-startByte)*8 + (endBit - startBit)
		if bits == 0 {
			continue
		}
		if bits > typeFieldBits {
			// bits = typeFieldBits
			return fmt.Errorf("The field is to small to store the num by tag:%s", field.Name)
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

		var n uint64
		for index := startByte; index <= endByte; index++ {
			lbit, rbit := 0, 8
			if index == startByte {
				lbit = startBit
			}
			if index == endByte {
				rbit = endBit
			}
			for bitIndex := lbit; bitIndex < rbit; bitIndex++ {
				n <<= 1
				n |= uint64((b[index] >> (8 - bitIndex - 1)) & 0x01)
			}
		}

		valueField := value.Field(i)
		if valueField.CanSet() {
			switch field.Type.Kind() {
			case reflect.Int:
				valueField.Set(reflect.ValueOf(int(n)))
			case reflect.Int8:
				valueField.Set(reflect.ValueOf(int8(n)))
			case reflect.Int16:
				valueField.Set(reflect.ValueOf(int16(n)))
			case reflect.Int32:
				valueField.Set(reflect.ValueOf(int32(n)))
			case reflect.Int64:
				valueField.Set(reflect.ValueOf(int64(n)))
			case reflect.Uint:
				valueField.Set(reflect.ValueOf(uint(n)))
			case reflect.Uint8:
				valueField.Set(reflect.ValueOf(uint8(n)))
			case reflect.Uint16:
				valueField.Set(reflect.ValueOf(uint16(n)))
			case reflect.Uint32:
				valueField.Set(reflect.ValueOf(uint32(n)))
			case reflect.Uint64:
				valueField.Set(reflect.ValueOf(uint64(n)))
			default:
				return fmt.Errorf("invalid field:%s", field.Name)
			}
		}
	}
	return nil
}

var (
	invalidErr = fmt.Errorf("empty expression")
)

//[startByte.startBit,endByte.endBit)
func parse(expression string) (startByte, startBit, endByte, endBit int, err error) {
	startByte, startBit, endByte, endBit = 0, 0, 0, 0
	if expression == "-" || expression == "" {
		return
	}

	if !strings.HasPrefix(expression, "[") || !strings.HasSuffix(expression, "]") {
		err = invalidErr
		return
	}

	var err1, err2, err3, err4 error
	slice := strings.Split(expression[1:len(expression)-1], ":")
	switch len(slice) {
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
				err = invalidErr
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
				err = invalidErr
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
				err = invalidErr
			}
		}
	default:
		err = invalidErr
	}

	err = easyerrors.HandleMultiError(easyerrors.Simple(), err1, err2, err3, err4)
	if startByte < -1 || endByte < -1 || startBit < 0 || startBit > 7 || endBit < 0 || endBit > 7 || endByte < startByte || (endByte == startByte && endBit < startBit) {
		err = invalidErr
	}
	return
}
