package easybits

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/SmartBrave/Athena/easyio"
)

//func TestMarshal(t *testing.T) {
//	type args struct {
//		v      interface{}
//		writer easyio.EasyWriter
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := Marshal(tt.args.v, tt.args.writer); (err != nil) != tt.wantErr {
//				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestUnmarshal(t *testing.T) {
	type foo struct {
		A uint8 `bits:"[1.2:2.2]"` //第一个字节的第二个比特，到第二个字节的第三个比特
		// B uint8 `bits:"[1]"`       //紧邻上面字段的后一个比特 TODO: support after
		// C uint8 `bits:"[3.1]"`     //第三个字节的第一个比特 TODO: support after
		D uint8 `bits:"-"` //ignore
		// E uint8 `bits:"[-]"` TODO: support after
		// F uint8 `bits:"[1]"` TODO: support after
		// G uint8 `bits:"[1]"` TODO: support after
		// H uint8 `bits:"[1]"` TODO: support after
		I uint8 `bits:"[2.3:2.4]"`
		J uint8 `bits:"[3.0:4.0]"`
		// TODO: `bits:[1.0:2.0],if xxx`
	}
	type args struct {
		b []byte
		v foo
	}
	tests := []struct {
		name    string
		args    args
		want    foo
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "first",
			args: args{
				b: []byte{0b10000011, 0b00010001, 0b00111001, 0b11111010},
				v: foo{},
			},
			want: foo{
				A: 0b01000100,
				// B: 0x00,
				// C: 0x00,
				D: 0b0,
				// E: 0x00,
				// F: 0x00,
				// G: 0x01,
				// H: 0x01,
				I: 0b1,
				J: 0b11111010,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesReader := bytes.NewReader(tt.args.b)
			reader := easyio.NewEasyReader(bytesReader)

			if err := Unmarshal(reader, &tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("Unmarshal fail. want:%v, got:%v", tt.want, tt.args.v)
			}
		})
	}
}

func TestUnmarshal_1(t *testing.T) {
	type foo struct {
		A     uint8    `bits:"[0.0:0.1]"`
		B     uint8    `bits:"[0.1:0.2]"`
		C     uint8    `bits:"[0.2:0.3]"`
		D     uint8    `bits:"[0.3:0.4]"`
		E     uint8    `bits:"[0.4:0.5]"`
		F     uint8    `bits:"[0.5:0.6]"`
		G     uint8    `bits:"[0.6:0.7]"`
		H     uint8    `bits:"[0.7:1.0]"`
		I     uint8    `bits:"[0.0:0.2]"`
		J     uint16   `bits:"[0.0:0.3]"`
		K     uint8    `bits:"[0.0:0.4]"`
		L     uint8    `bits:"[0.0:0.5]"`
		M     uint8    `bits:"[0.0:0.6]"`
		N     uint8    `bits:"[0.0:0.7]"`
		O     uint8    `bits:"[0.0:1.0]"`
		P     uint16   `bits:"[0.0:1.1]"`
		Q     uint16   `bits:"[0.0:1.2]"`
		R     uint16   `bits:"[0.0:1.3]"`
		S     uint16   `bits:"[0.0:1.4]"`
		T     uint16   `bits:"[0.0:1.5]"`
		U     uint16   `bits:"[0.0:1.6]"`
		V     uint16   `bits:"[0.0:1.7]"`
		W     uint8    `bits:"[0.4:1.1]"`
		X     uint16   `bits:"[0.4:1.7]"`
		Y     uint16   `bits:"[0.7:2.0]"`
		Z     uint16   `bits:"[1.3:2.7]"`
		Array [7]uint8 `bits:"[0.0:2.5:3]"`
		Slice []uint8  `bits:"[0.0:2.6:5]"`
	}
	type args struct {
		b []byte
		v foo
	}
	tests := []struct {
		name    string
		args    args
		want    foo
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "first",
			args: args{
				b: []byte{0b10101010, 0b11011011, 0b00100100},
				v: foo{},
			},

			want: foo{
				A: 0b1,
				B: 0b0,
				C: 0b1,
				D: 0b0,
				E: 0b1,
				F: 0b0,
				G: 0b1,
				H: 0b0,
				I: 0b10,
				J: 0b101,
				K: 0b1010,
				L: 0b10101,
				M: 0b101010,
				N: 0b1010101,
				O: 0b10101010,
				P: 0b101010101,
				Q: 0b1010101011,
				R: 0b10101010110,
				S: 0b101010101101,
				T: 0b1010101011011,
				U: 0b10101010110110,
				V: 0b101010101101101,
				W: 0b10101,
				X: 0b10101101101,
				Y: 0b011011011,
				Z: 0b110110010010,
				Array: [7]uint8{
					0b101,
					0b010,
					0b101,
					0b101,
					0b101,
					0b100,
					0b100,
				},
				Slice: []uint8{
					0b10101,
					0b01011,
					0b01101,
					0b10010,
					0b01,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesReader := bytes.NewReader(tt.args.b)
			reader := easyio.NewEasyReader(bytesReader)

			if err := Unmarshal(reader, &tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("Unmarshal fail. want:%v, got:%v", tt.want, tt.args.v)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		expression    string
		wantStartByte int
		wantStartBit  int
		wantEndByte   int
		wantEndBit    int
		wantItemBit   int
		wantErr       bool
	}{
		{
			name:          "first",
			expression:    "",
			wantStartByte: 0,
			wantStartBit:  0,
			wantEndByte:   0,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "second",
			expression:    "-",
			wantStartByte: 0,
			wantStartBit:  0,
			wantEndByte:   0,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "third",
			expression:    "[1.0:2.0]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "fourth",
			expression:    "[1:2]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "fifth",
			expression:    "[1:2.1]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    1,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "sixth",
			expression:    "[1.1:2]",
			wantStartByte: 1,
			wantStartBit:  1,
			wantEndByte:   2,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "seventh",
			expression:    "[:2]",
			wantStartByte: 0,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "eighth",
			expression:    "[:2.1]", //TODO: support
			wantStartByte: 0,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    1,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "ninth",
			expression:    "[1:]", //TODO: support
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   0,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       true,
		},
		{
			name:          "tenth",
			expression:    "[1.1:]", //TODO: support
			wantStartByte: 1,
			wantStartBit:  1,
			wantEndByte:   0,
			wantEndBit:    0,
			wantItemBit:   0,
			wantErr:       true,
		},
		{
			name:          "eleventh",
			expression:    "[1.1]", //TODO: support
			wantStartByte: 1,
			wantStartBit:  1,
			wantEndByte:   1,
			wantEndBit:    2,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:          "twelfth",
			expression:    "[1]", //TODO: support
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   1,
			wantEndBit:    7,
			wantItemBit:   0,
			wantErr:       false,
		},
		{
			name:       "thirteenth",
			expression: "[1.:2]",
			//wantStartByte: 1,
			//wantStartBit:  0,
			//wantEndByte:   1,
			//wantEndBit:    7,
			// wantItemBit:0,
			wantErr: true,
		},
		{
			name:          "fourteenth",
			expression:    "[0.0:3.5:3]",
			wantStartByte: 0,
			wantStartBit:  0,
			wantEndByte:   3,
			wantEndBit:    5,
			wantItemBit:   3,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if startByte, startBit, endByte, endBit, itemBit, err := parse(tt.expression); (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && (startByte != tt.wantStartByte || startBit != tt.wantStartBit || endByte != tt.wantEndByte || endBit != tt.wantEndBit || itemBit != tt.wantItemBit) {
				t.Errorf("parse() startByte = %d, startBit = %d, endByte = %d, endBit = %d, wantStartByte = %d, wantStartBit = %d, wantEndByte = %d, wantEndBit = %d, wantItemBit = %d", startByte, startBit, endByte, endBit, tt.wantStartByte, tt.wantStartBit, tt.wantEndByte, tt.wantEndBit, tt.wantItemBit)
			}
		})
	}
}
