package bitmap

import (
	"testing"
)

// Bitmap is 1-indexed: position 1 is the lowest bit of the first byte.
func TestSetGetRoundTrip(t *testing.T) {
	bm, err := New(20)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []uint32{1, 2, 8, 9, 16, 20} {
		if err := bm.Set(p); err != nil {
			t.Fatalf("Set(%d): %v", p, err)
		}
	}
	for _, p := range []uint32{1, 2, 8, 9, 16, 20} {
		if !bm.Get(p) {
			t.Errorf("Get(%d)=false, want true", p)
		}
	}
	for _, p := range []uint32{3, 4, 5, 6, 7, 10, 11, 17, 18, 19} {
		if bm.Get(p) {
			t.Errorf("Get(%d)=true on unset bit", p)
		}
	}
}

func TestSetOutOfRange(t *testing.T) {
	bm, _ := New(8)
	if err := bm.Set(0); err != ERROR {
		t.Errorf("Set(0) err = %v, want ERROR (1-indexed)", err)
	}
	if err := bm.Set(1000); err != ERROR {
		t.Errorf("Set(1000) err = %v, want ERROR", err)
	}
}

func TestResetClearsBit(t *testing.T) {
	bm, _ := New(8)
	bm.Set(3)
	if !bm.Get(3) {
		t.Fatal("precondition failed: Set(3) didn't take")
	}
	if err := bm.Reset(3); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if bm.Get(3) {
		t.Errorf("Get(3) after Reset = true")
	}
}

func TestResetOutOfRange(t *testing.T) {
	bm, _ := New(8)
	if err := bm.Reset(0); err != ERROR {
		t.Errorf("Reset(0) err = %v, want ERROR", err)
	}
	if err := bm.Reset(99); err != ERROR {
		t.Errorf("Reset(99) err = %v, want ERROR", err)
	}
}

func TestGetOutOfRangeReturnsFalse(t *testing.T) {
	bm, _ := New(8)
	if bm.Get(0) {
		t.Errorf("Get(0) = true on 1-indexed bitmap, want false")
	}
	if bm.Get(99) {
		t.Errorf("Get(99) = true past end, want false")
	}
}

func TestRangeVisitsAllPositions(t *testing.T) {
	bm, _ := New(10)
	visited := make([]uint32, 0, 10)
	bm.Range(func(pos uint32) { visited = append(visited, pos) })
	if len(visited) != 10 {
		t.Fatalf("Range visited %d positions, want 10", len(visited))
	}
	for i, p := range visited {
		if p != uint32(i+1) {
			t.Errorf("Range[%d]=%d, want %d (1-indexed)", i, p, i+1)
		}
	}
}

func TestNewZeroBitsErrors(t *testing.T) {
	bm, err := New(0)
	if err != ERROR || bm != nil {
		t.Errorf("New(0) = (%v, %v), want (nil, ERROR)", bm, err)
	}
}

func TestNewWithStringRoundTrip(t *testing.T) {
	src := "abc"
	bm, err := NewWithString(src)
	if err != nil {
		t.Fatalf("NewWithString: %v", err)
	}
	if got := bm.String(); got != src {
		t.Errorf("NewWithString().String() = %q, want %q", got, src)
	}
	// Each byte from `src` becomes 8 bits in the bitmap, so individual bits
	// should match source bytes.
	for i, c := range []byte(src) {
		for bit := 0; bit < 8; bit++ {
			pos := uint32(i*BYTE + bit + 1)
			expected := (c>>bit)&0x1 == 1
			if got := bm.Get(pos); got != expected {
				t.Errorf("Get(%d)=%v, want %v (byte=%#08b bit=%d)", pos, got, expected, c, bit)
			}
		}
	}
}

func TestStringReflectsBuffer(t *testing.T) {
	bm, _ := New(16)
	bm.Set(1)
	bm.Set(2)
	// bits 1+2 set in first byte = 0b00000011 = 0x03
	if got := bm.String(); len(got) != 2 || got[0] != 0x03 || got[1] != 0x00 {
		t.Errorf("String() = %x, want first byte = 0x03", got)
	}
}
