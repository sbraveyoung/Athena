package easysyntax

import (
	"errors"
	"testing"
)

func TestItoa_64(t *testing.T) {
	tests := []struct {
		name string
		in   uint64
		want string
	}{
		{"zero", 0, "0"},
		{"nine", 9, "9"},
		{"ten", 10, "a"},
		{"thirty_five", 35, "z"},
		{"thirty_six", 36, "A"},
		{"sixty_one", 61, "Z"},
		{"sixty_two", 62, "."},
		{"sixty_three", 63, "-"},
		{"sixty_four", 64, "10"},
		{"composed", (((((1*64+10)*64+2)*64+37)*64+3)*64+62)*64 + 63, "1a2B3.-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Itoa_64(tt.in); got != tt.want {
				t.Errorf("Itoa_64(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestAtoi_64(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    uint64
		wantErr error
	}{
		{"empty", "", 0, nil},
		{"zero", "0", 0, nil},
		{"nine", "9", 9, nil},
		{"a", "a", 10, nil},
		{"z", "z", 35, nil},
		{"A", "A", 36, nil},
		{"Z", "Z", 61, nil},
		{"dot", ".", 62, nil},
		{"dash", "-", 63, nil},
		{"two_digit", "10", 64, nil},
		{"composed", "1a2B3.-", (((((1*64+10)*64+2)*64+37)*64+3)*64+62)*64 + 63, nil},
		{"invalid_plus", "+", 0, ErrSyntax},
		{"invalid_slash", "/", 0, ErrSyntax},
		{"invalid_mid", "1a/", 0, ErrSyntax},
		{"invalid_space", " ", 0, ErrSyntax},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Atoi_64(tt.in)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Atoi_64(%q) unexpected err: %v", tt.in, err)
				}
				if got != tt.want {
					t.Errorf("Atoi_64(%q) = %d, want %d", tt.in, got, tt.want)
				}
				return
			}
			if err == nil {
				t.Fatalf("Atoi_64(%q) expected err %v, got nil", tt.in, tt.wantErr)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Atoi_64(%q) err = %v, want errors.Is == %v", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestAtoiItoaRoundTrip(t *testing.T) {
	values := []uint64{0, 1, 63, 64, 4095, 4096, 1 << 30, (1 << 60) - 1}
	for _, v := range values {
		s := Itoa_64(v)
		got, err := Atoi_64(s)
		if err != nil {
			t.Errorf("round-trip %d -> %q decode err: %v", v, s, err)
			continue
		}
		if got != v {
			t.Errorf("round-trip %d -> %q -> %d", v, s, got)
		}
	}
}

func TestNumErrorUnwrap(t *testing.T) {
	_, err := Atoi_64("+")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrSyntax) {
		t.Errorf("errors.Is(err, ErrSyntax) = false, err = %v", err)
	}
	var ne *NumError
	if !errors.As(err, &ne) {
		t.Errorf("errors.As(err, *NumError) = false")
	}
	if ne.Func != "Atoi_64" || ne.Num != "+" {
		t.Errorf("NumError fields: Func=%q Num=%q", ne.Func, ne.Num)
	}
	if got := ne.Error(); got == "" {
		t.Errorf("NumError.Error() empty")
	}
}
