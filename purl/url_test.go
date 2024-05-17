package purl

import (
	"net/url"
	"reflect"
	"testing"
)

var ki = &indexTable{}

func init() {
	ki.forwardIndex.Store("a", 1)
	ki.forwardIndex.Store("b", 2)
	ki.forwardIndex.Store("c", 3)
	ki.forwardIndex.Store("d", 4)
	ki.forwardIndex.Store("e", 5)
	ki.forwardIndex.Store("f", 6)
	ki.forwardIndex.Store("g", 7)
	ki.forwardIndex.Store("h", 8)
	ki.forwardIndex.Store("i", 9)
	ki.forwardIndex.Store("j", 10)
	ki.forwardIndex.Store("k", 11)
	ki.forwardIndex.Store("l", 12)
	ki.forwardIndex.Store("m", 13)
	ki.forwardIndex.Store("n", 14)
	ki.forwardIndex.Store("o", 15)
	ki.forwardIndex.Store("p", 16)
	ki.forwardIndex.Store("q", 17)
	ki.forwardIndex.Store("r", 18)
	ki.forwardIndex.Store("s", 19)
	ki.forwardIndex.Store("t", 20)
	ki.forwardIndex.Store("u", 21)
	ki.forwardIndex.Store("v", 22)
	ki.forwardIndex.Store("w", 23)
	ki.forwardIndex.Store("x", 24)
	ki.forwardIndex.Store("y", 25)
	ki.forwardIndex.Store("z", 26)

	ki.backwardIndex.Store(1, "a")
	ki.backwardIndex.Store(2, "b")
	ki.backwardIndex.Store(3, "c")
	ki.backwardIndex.Store(4, "d")
	ki.backwardIndex.Store(5, "e")
	ki.backwardIndex.Store(6, "f")
	ki.backwardIndex.Store(7, "g")
	ki.backwardIndex.Store(8, "h")
	ki.backwardIndex.Store(9, "i")
	ki.backwardIndex.Store(10, "j")
	ki.backwardIndex.Store(11, "k")
	ki.backwardIndex.Store(12, "l")
	ki.backwardIndex.Store(13, "m")
	ki.backwardIndex.Store(14, "n")
	ki.backwardIndex.Store(15, "o")
	ki.backwardIndex.Store(16, "p")
	ki.backwardIndex.Store(17, "q")
	ki.backwardIndex.Store(18, "r")
	ki.backwardIndex.Store(19, "s")
	ki.backwardIndex.Store(20, "t")
	ki.backwardIndex.Store(21, "u")
	ki.backwardIndex.Store(22, "v")
	ki.backwardIndex.Store(23, "w")
	ki.backwardIndex.Store(24, "x")
	ki.backwardIndex.Store(25, "y")
	ki.backwardIndex.Store(26, "z")
}

func TestParseQuery(t *testing.T) {
	type args struct {
		rawQuery string
	}
	tests := []struct {
		name    string
		args    args
		wantM   Values
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "first",
			args: args{
				rawQuery: "a=b&c=d&e=f",
			},
			wantM: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
					5: 6,
				},
				sv: map[string]string{},
			},
		},
		{
			name: "second",
			args: args{
				rawQuery: "a=b",
			},
			wantM: Values{
				iv: map[int]int{
					1: 2,
				},
				sv: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := ParseQuery(tt.args.rawQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery return err = %v, want err = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotM, tt.wantM) {
				t.Errorf("ParseQuery() = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}

func TestValues_Add(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		got  Values
		args args
		want Values
	}{
		// TODO: Add test cases.
		{
			name: "first",
			got: Values{
				iv: map[int]int{},
				sv: map[string]string{},
			},
			args: args{
				key:   "a",
				value: "b",
			},
			want: Values{
				iv: map[int]int{
					1: 2,
				},
				sv: map[string]string{},
			},
		},
		{
			name: "second",
			got: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
				},
				sv: map[string]string{},
			},
			args: args{
				key:   "e",
				value: "f",
			},
			want: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
					5: 6,
				},
				sv: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.Add(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Values_Add() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestValues_Del(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		got  Values
		args args
		want Values
	}{
		// TODO: Add test cases.
		{
			name: "first",
			got: Values{
				iv: map[int]int{
					1: 2,
				},
				sv: map[string]string{},
			},
			args: args{
				key: "a",
			},
			want: Values{
				iv: map[int]int{},
				sv: map[string]string{},
			},
		},
		{
			name: "second",
			got: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
				},
				sv: map[string]string{},
			},
			args: args{
				key: "c",
			},
			want: Values{
				iv: map[int]int{
					1: 2,
				},
				sv: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.Del(tt.args.key)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Values_Del() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

func TestValues_Encode(t *testing.T) {
	tests := []struct {
		name string
		v    Values
		want string
	}{
		// TODO: Add test cases.
		{
			name: "first",
			v: Values{
				iv: map[int]int{
					1: 2,
				},
				sv: map[string]string{},
			},
			want: "a=b",
		},
		{
			name: "second",
			v: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
				},
				sv: map[string]string{},
			},
			want: "a=b&c=d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.Encode()
			if got != tt.want {
				t.Errorf("Values.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		v    Values
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "first",
			v: Values{
				iv: map[int]int{
					1: 2,
				},
				sv: map[string]string{},
			},
			args: args{
				key: "a",
			},
			want: "b",
		},
		{
			name: "second",
			v: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
					5: 6,
				},
				sv: map[string]string{},
			},
			args: args{
				key: "c",
			},
			want: "d",
		},
		{
			name: "third",
			v: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
					5: 6,
				},
				sv: map[string]string{},
			},
			args: args{
				key: "g",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Values.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		got  Values
		args args
		want Values
	}{
		// TODO: Add test cases.
		{
			name: "first",
			got: Values{
				iv: map[int]int{},
				sv: map[string]string{},
			},
			args: args{
				key:   "e",
				value: "f",
			},
			want: Values{
				iv: map[int]int{
					5: 6,
				},
				sv: map[string]string{},
			},
		},
		{
			name: "second",
			got: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
				},
				sv: map[string]string{},
			},
			args: args{
				key:   "e",
				value: "f",
			},
			want: Values{
				iv: map[int]int{
					1: 2,
					3: 4,
					5: 6,
				},
				sv: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.got.Set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Values_Set() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

// ------------------------------Benchmark---------------------------------
var (
	// rawQuery = "a=1&b=2&c=3&d=4&e=5&f=6&g=7&h=8&i=9&j=10&k=11&l=12&m=13&n=14&o=15&p=16&q=17&r=18&s=19&t=20&u=21&v=22&w=23&x=24&y=25&z=26"
	// rawQuery = "a=1"
	rawQuery = "codecInfo=8192&dpSrc=13&ikChorus=1&ikDebug=1&ikDnsOp=1001&ikFastRate=1.1&ikHost=ik&ikLog=1&ikMaxBuf=3600&ikMinBuf=2900&ikOp=0&ikSlowRate=0.9&ikSyncBetaxyz=0&msSmid=D2nD1k5JCYy1r7c9Bs4KRxYIQyJtTLaUtLUG7BadHTvZ4X4d&msUid=717620689&pushHost=push.cls.inke.cn"
)

func BenchmarkParseQuery(b *testing.B) {
	var v Values
	var s string
	for n := 0; n < b.N; n++ {
		// XXX: 性能差查一下 go 版本的问题
		v, _ = ParseQuery(rawQuery)
		s = v.Encode()
	}
	_ = s
	_ = v
}

func BenchmarkStdParseQuery(b *testing.B) {
	var v url.Values
	var s string
	for n := 0; n < b.N; n++ {
		v, _ = url.ParseQuery(rawQuery)
		s = v.Encode()
	}
	_ = s
	_ = v
}
