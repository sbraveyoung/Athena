package bitmap

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		bits uint32
	}
	type want struct {
		bm  *Bitmap
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
		{
			name: "first",
			args: args{
				bits: 0,
			},
			want: want{
				bm:  nil,
				err: ERROR,
			},
		},
		{
			name: "second",
			args: args{
				bits: 1,
			},
			want: want{
				bm: &Bitmap{
					bits:   1,
					buffer: []byte{0x0},
				},
				err: nil,
			},
		},
		{
			name: "third",
			args: args{
				bits: 7,
			},
			want: want{
				bm: &Bitmap{
					bits:   7,
					buffer: []byte{0x0},
				},
				err: nil,
			},
		},
		{
			name: "fourth",
			args: args{
				bits: 8,
			},
			want: want{
				bm: &Bitmap{
					bits:   8,
					buffer: []byte{0x0},
				},
				err: nil,
			},
		},
		{
			name: "fifth",
			args: args{
				bits: 9,
			},
			want: want{
				bm: &Bitmap{
					bits:   9,
					buffer: []byte{0x0, 0x0},
				},
				err: nil,
			},
		},
		{
			name: "sixth",
			args: args{
				bits: 16,
			},
			want: want{
				bm: &Bitmap{
					bits:   16,
					buffer: []byte{0x0, 0x0},
				},
				err: nil,
			},
		},
		{
			name: "seventh",
			args: args{
				bits: 20,
			},
			want: want{
				bm: &Bitmap{
					bits:   20,
					buffer: []byte{0x0, 0x0, 0x0},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(tt.args.bits); !reflect.DeepEqual(got, tt.want.bm) || err != tt.want.err {
				t.Errorf("New(%d) = %v, %v, want %v", tt.args.bits, got, err, tt.want)
			}
		})
	}
}

func TestBitmap_Set(t *testing.T) {
	type fields struct {
		buffer []byte
	}
	type args struct {
		pos uint32
	}
	type want struct {
		bm  *Bitmap
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		// TODO: Add test cases.
		{
			name: "first",
			fields: fields{
				buffer: []byte{0x0},
			},
			args: args{
				pos: 0,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x0},
				},
				err: ERROR,
			},
		},
		{
			name: "second",
			fields: fields{
				buffer: []byte{0x0, 0x0},
			},
			args: args{
				pos: 1,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x1, 0x0},
				},
				err: nil,
			},
		},
		{
			name: "third",
			fields: fields{
				buffer: []byte{0x3, 0x5},
			},
			args: args{
				pos: 8,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x83, 0x5},
				},
				err: nil,
			},
		},
		{
			name: "fourth",
			fields: fields{
				buffer: []byte{0xff, 0x8},
			},
			args: args{
				pos: 10,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0xff, 0xa},
				},
				err: nil,
			},
		},
		{
			name: "fifth",
			fields: fields{
				buffer: []byte{0x0, 0xff, 0xf3},
			},
			args: args{
				pos: 20,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x0, 0xff, 0xfb},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := &Bitmap{
				buffer: tt.fields.buffer,
			}
			if err := bm.Set(tt.args.pos); !reflect.DeepEqual(bm, tt.want.bm) || err != tt.want.err {
				t.Errorf("bm.Set(%d) = %v, %v, want %v", tt.args.pos, bm, err, tt.want)
			}
		})
	}
}

func TestBitmap_Reset(t *testing.T) {
	type fields struct {
		buffer []byte
	}
	type args struct {
		pos uint32
	}
	type want struct {
		bm  *Bitmap
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		// TODO: Add test cases.
		{
			name: "first",
			fields: fields{
				buffer: []byte{0x0},
			},
			args: args{
				pos: 0,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x0},
				},
				err: ERROR,
			},
		},
		{
			name: "second",
			fields: fields{
				buffer: []byte{0xff, 0xff},
			},
			args: args{
				pos: 1,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0xfe, 0xff},
				},
				err: nil,
			},
		},
		{
			name: "third",
			fields: fields{
				buffer: []byte{0x3, 0x5},
			},
			args: args{
				pos: 8,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x03, 0x5},
				},
				err: nil,
			},
		},
		{
			name: "fourth",
			fields: fields{
				buffer: []byte{0xff, 0x8},
			},
			args: args{
				pos: 8,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x7f, 0x8},
				},
				err: nil,
			},
		},
		{
			name: "fifth",
			fields: fields{
				buffer: []byte{0x0, 0xff, 0xf3},
			},
			args: args{
				pos: 18,
			},
			want: want{
				bm: &Bitmap{
					buffer: []byte{0x0, 0xff, 0xf1},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := &Bitmap{
				buffer: tt.fields.buffer,
			}
			if err := bm.Reset(tt.args.pos); !reflect.DeepEqual(bm, tt.want.bm) || err != tt.want.err {
				t.Errorf("bm.Reset(%d) = %v, %v, want %v", tt.args.pos, bm, err, tt.want)
			}
		})
	}
}

func TestBitmap_Get(t *testing.T) {
	type fields struct {
		buffer []byte
	}
	type args struct {
		pos uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "first",
			fields: fields{
				buffer: []byte{0x17, 0x85, 0x26},
			},
			args: args{
				pos: 0,
			},
			want: false,
		},
		{
			name: "second",
			fields: fields{
				buffer: []byte{0x17, 0x85, 0x26},
			},
			args: args{
				pos: 1,
			},
			want: true,
		},
		{
			name: "third",
			fields: fields{
				buffer: []byte{0x17, 0x85, 0x26},
			},
			args: args{
				pos: 7,
			},
			want: false,
		},
		{
			name: "fourth",
			fields: fields{
				buffer: []byte{0x17, 0x85, 0x26},
			},
			args: args{
				pos: 8,
			},
			want: false,
		},
		{
			name: "fifth",
			fields: fields{
				buffer: []byte{0x17, 0x85, 0x26},
			},
			args: args{
				pos: 30,
			},
			want: false,
		},
		{
			name: "sixth",
			fields: fields{
				buffer: []byte{0x17, 0x85, 0x26},
			},
			args: args{
				pos: 22,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm := &Bitmap{
				buffer: tt.fields.buffer,
			}
			if got := bm.Get(tt.args.pos); got != tt.want {
				t.Errorf("Bitmap.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
