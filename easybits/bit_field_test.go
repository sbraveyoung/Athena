package easybits

import (
	"testing"
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
//
//func TestUnmarshal(t *testing.T) {
//	type foo struct {
//		A uint8 `bits:"[1.2:2.3]"` //第一个字节的第二个比特，到第二个字节的第三个比特
//		B uint8 `bits:"[1]"`       //紧邻上面字段的后一个比特
//		C uint8 `bits:"[3.1]"`     //第三个字节的第一个比特
//		D uint8 `bits:"-"`         //ignore
//		E uint8 `bits:"[-]"`
//		F uint8 `bits:"[1]"`
//		G uint8 `bits:"[1]"`
//		H uint8 `bits:"[1]"`
//	}
//	type args struct {
//		b []byte
//		v foo
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    foo
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//		{
//			name: "first",
//			args: args{
//				b: []byte{0x83},
//				v: foo{},
//			},
//			want: foo{
//				A: 0x01,
//				B: 0x00,
//				C: 0x00,
//				D: 0x00,
//				E: 0x00,
//				F: 0x00,
//				G: 0x01,
//				H: 0x01,
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			bytesReader := bytes.NewReader(tt.args.b)
//			reader := easyio.NewEasyReader(bytesReader)
//
//			if err := Unmarshal(reader, &tt.args.v); (err != nil) != tt.wantErr {
//				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
//			}
//			if !reflect.DeepEqual(tt.args.v, tt.want) {
//				t.Errorf("Unmarshal fail. want:%v, got:%v", tt.want, tt.args.v)
//			}
//		})
//	}
//}

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		expression    string
		wantStartByte int
		wantStartBit  int
		wantEndByte   int
		wantEndBit    int
		wantErr       bool
	}{
		{
			name:          "first",
			expression:    "",
			wantStartByte: -1,
			wantStartBit:  -1,
			wantEndByte:   -1,
			wantEndBit:    -1,
			wantErr:       false,
		},
		{
			name:          "second",
			expression:    "-",
			wantStartByte: -1,
			wantStartBit:  -1,
			wantEndByte:   -1,
			wantEndBit:    -1,
			wantErr:       false,
		},
		{
			name:          "third",
			expression:    "[1.0:2.0]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    0,
			wantErr:       false,
		},
		{
			name:          "fourth",
			expression:    "[1:2]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    0,
			wantErr:       false,
		},
		{
			name:          "fifth",
			expression:    "[1:2.1]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    1,
			wantErr:       false,
		},
		{
			name:          "sixth",
			expression:    "[1.1:2]",
			wantStartByte: 1,
			wantStartBit:  1,
			wantEndByte:   2,
			wantEndBit:    0,
			wantErr:       false,
		},
		{
			name:          "seventh",
			expression:    "[:2]",
			wantStartByte: -1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    0,
			wantErr:       false,
		},
		{
			name:          "eighth",
			expression:    "[:2.1]",
			wantStartByte: -1,
			wantStartBit:  0,
			wantEndByte:   2,
			wantEndBit:    1,
			wantErr:       false,
		},
		{
			name:          "ninth",
			expression:    "[1:]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   -1,
			wantEndBit:    0,
			wantErr:       false,
		},
		{
			name:          "tenth",
			expression:    "[1.1:]",
			wantStartByte: 1,
			wantStartBit:  1,
			wantEndByte:   -1,
			wantEndBit:    0,
			wantErr:       false,
		},
		{
			name:          "eleventh",
			expression:    "[1.1]",
			wantStartByte: 1,
			wantStartBit:  1,
			wantEndByte:   1,
			wantEndBit:    2,
			wantErr:       false,
		},
		{
			name:          "twelfth",
			expression:    "[1]",
			wantStartByte: 1,
			wantStartBit:  0,
			wantEndByte:   1,
			wantEndBit:    7,
			wantErr:       false,
		},
		{
			name:       "thirteenth",
			expression: "[1.:2]",
			//wantStartByte: 1,
			//wantStartBit:  0,
			//wantEndByte:   1,
			//wantEndBit:    7,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if startByte, startBit, endByte, endBit, err := parse(tt.expression); (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil && (startByte != tt.wantStartByte || startBit != tt.wantStartBit || endByte != tt.wantEndByte || endBit != tt.wantEndBit) {
				t.Errorf("parse() startByte = %d, startBit = %d, endByte = %d, endBit = %d, wantStartByte = %d, wantStartBit = %d, wantEndByte = %d, wantEndBit = %d", startByte, startBit, endByte, endBit, tt.wantStartByte, tt.wantStartBit, tt.wantEndByte, tt.wantEndBit)
			}
		})
	}
}
