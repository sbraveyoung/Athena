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
	return nil
}

func Unmarshal(reader easyio.EasyReader, v interface{}) (err error) {
	//value := reflect.ValueOf(v)
	//if value.Kind() != reflect.Ptr {
	//	return fmt.Errorf("invalid type of v:%+v", value.Kind())
	//}
	//value = value.Elem()
	//if value.Kind() != reflect.Struct {
	//	return fmt.Errorf("invalid type of *v:%+v", value.Kind())
	//}

	tYpe := reflect.TypeOf(v)
	if tYpe.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid type of v:%+v", tYpe.Kind())
	}
	tYpe = tYpe.Elem()
	if tYpe.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid type of v:%+v", tYpe.Kind())
	}

	//var b []byte
	//var index int
	//for i := 0; i < tYpe.NumField(); i++ {
	//	field := tYpe.Field(i)
	//	bitTag := field.Tag.Get("bits")

	//	bits := field.Type.Bits()
	//	if bitTag != "" {
	//		bits, err = strconv.Atoi(bitTag)
	//		if err != nil {
	//			return fmt.Errorf("invlid tag, %v", err)
	//		}
	//		if bits > typeFieldBits {
	//			bits = typeFieldBits
	//		}
	//	}

	//	if len(b)*8-index < bits {
	//		needBits := bits - (len(b)*8 - index)
	//		needBytes := needBits / 8
	//		if needBits%8 != 0 {
	//			needBytes++
	//		}
	//		b = append(b, make([]byte, needBytes)...)
	//		err = reader.ReadFull(b[len(b)-needBytes:])
	//		if err != nil {
	//			return err
	//		}
	//	}

	//	valueField := value.Field(i)
	//	if 8-index >= bits {
	//		// valueField.Set(b[0]>>index
	//	}
	//}
	return nil
}

var (
	invalidErr = fmt.Errorf("empty expression")
)

//[startByte.startBit,endByte.endBit)
//startByte == -1 means current byte
//endByte == -1 means all data to fill this field
//startBit == -1 or endBit == -1 means ignore this field
func parse(expression string) (startByte, startBit, endByte, endBit int, err error) {
	if expression == "-" || expression == "" {
		return -1, -1, -1, -1, nil
	}

	if !strings.HasPrefix(expression, "[") || !strings.HasSuffix(expression, "]") {
		err = invalidErr
		return
	}

	startByte, startBit, endByte, endBit = -1, 0, -1, 0
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
	if startByte < -1 || endByte < -1 || startBit < 0 || startBit > 7 || endBit < 0 || endBit > 7 {
		err = invalidErr
	}
	return
}
