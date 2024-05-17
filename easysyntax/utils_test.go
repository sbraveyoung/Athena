package easysyntax

import (
	"testing"
)

func TestItoa_64(t *testing.T) {
	type args struct {
		i uint64
	}
	tests := []struct {
		name  string
		args  args
		wantS string
	}{
		// TODO: Add test cases.
		{
			name: "0",
			args: args{
				i: 0,
			},
			wantS: "0",
		},
		{
			name: "10",
			args: args{
				i: 10,
			},
			wantS: "a",
		},
		{
			name: "36",
			args: args{
				i: 36,
			},
			wantS: "A",
		},
		{
			name: "63",
			args: args{
				i: 63,
			},
			wantS: "/",
		},
		{
			name: "64",
			args: args{
				i: 64,
			},
			wantS: "10",
		},
		{
			name: "79,500,165,055",
			args: args{
				i: (((((1*64+10)*64+2)*64+37)*64+3)*64+62)*64 + 63,
			},
			wantS: "1a2B3+/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := Itoa_64(tt.args.i); gotS != tt.wantS {
				t.Errorf("tt.name:%s, Itoa_64() = %v, want %v", tt.name, gotS, tt.wantS)
			}
		})
	}
}

func TestAtoi_64(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantN   uint64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "0",
			args: args{
				s: "0",
			},
			wantN:   0,
			wantErr: false,
		},
		{
			name: "10",
			args: args{
				s: "a",
			},
			wantN:   10,
			wantErr: false,
		},
		{
			name: "36",
			args: args{
				s: "A",
			},
			wantN:   36,
			wantErr: false,
		},
		{
			name: "63",
			args: args{
				s: "/",
			},
			wantN:   63,
			wantErr: false,
		},
		{
			name: "64",
			args: args{
				s: "10",
			},
			wantN:   64,
			wantErr: false,
		},
		{
			name: "79,500,165,055",
			args: args{
				s: "1a2B3+/",
			},
			wantN:   (((((1*64+10)*64+2)*64+37)*64+3)*64+62)*64 + 63,
			wantErr: false,
		},
		{
			name: "0",
			args: args{
				s: "0",
			},
			wantN:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := Atoi_64(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Atoi_64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Atoi_64() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
