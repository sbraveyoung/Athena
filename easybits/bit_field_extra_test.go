package easybits

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/SmartBrave/Athena/easyerrors"
	"github.com/SmartBrave/Athena/easyio"
)

func TestUnmarshalRejectsNonPointer(t *testing.T) {
	type s struct {
		A uint8 `bits:"[0.0:0.4]"`
	}
	r := easyio.NewEasyReader(bytes.NewReader([]byte{0xFF}))
	if err := Unmarshal(r, s{}); err == nil {
		t.Error("Unmarshal(non-pointer) returned nil err")
	}
}

func TestUnmarshalRejectsPointerToNonStruct(t *testing.T) {
	var n int
	r := easyio.NewEasyReader(bytes.NewReader([]byte{0xFF}))
	if err := Unmarshal(r, &n); err == nil {
		t.Error("Unmarshal(*int) returned nil err")
	}
}

func TestUnmarshalErrorOnReaderShortage(t *testing.T) {
	type s struct {
		A uint16 `bits:"[0.0:2.0]"` // wants 2 full bytes
	}
	r := easyio.NewEasyReader(bytes.NewReader([]byte{0xAB})) // only 1 byte available
	var got s
	if err := Unmarshal(r, &got); err == nil {
		t.Error("Unmarshal expected error on short reader, got nil")
	}
}

func TestUnmarshalRejectsInvalidTag(t *testing.T) {
	type s struct {
		A uint8 `bits:"not-a-tag"`
	}
	r := easyio.NewEasyReader(bytes.NewReader([]byte{0xFF}))
	if err := Unmarshal(r, &s{}); err == nil {
		t.Error("Unmarshal expected error on invalid tag, got nil")
	}
}

func TestUnmarshalRejectsFieldTooSmall(t *testing.T) {
	type s struct {
		A uint8 `bits:"[0.0:2.0]"` // 16 bits into a uint8 field
	}
	r := easyio.NewEasyReader(bytes.NewReader([]byte{0xFF, 0xFF}))
	if err := Unmarshal(r, &s{}); err == nil {
		t.Error("Unmarshal expected error: 16 bits cannot fit in uint8")
	}
}

func TestUnmarshalSkipsUnexported(t *testing.T) {
	// Unexported fields must be silently skipped, not panic.
	type s struct {
		A uint8 `bits:"[0.0:0.4]"`
		// b is unexported, even with a tag it should be ignored.
		b uint8 `bits:"[0.4:1.0]"` //nolint:unused
	}
	r := easyio.NewEasyReader(bytes.NewReader([]byte{0xFF}))
	var got s
	if err := Unmarshal(r, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.A != 0xF {
		t.Errorf("A=%x want 0xF", got.A)
	}
	_ = got.b
}

func TestParseEmptyTagAndDash(t *testing.T) {
	for _, in := range []string{"", "-"} {
		startB, startBit, endB, endBit, item, err := parse(in)
		if err != nil {
			t.Errorf("parse(%q): unexpected err %v", in, err)
		}
		if startB != 0 || startBit != 0 || endB != 0 || endBit != 0 || item != 0 {
			t.Errorf("parse(%q) returned non-zero", in)
		}
	}
}

func TestParseInvalidExpressions(t *testing.T) {
	bad := []string{
		"[",
		"]",
		"foo",
		"[1.:2]",
		"[1.x:2]",
		"[1:2:3:4]",
	}
	for _, in := range bad {
		if _, _, _, _, _, err := parse(in); err == nil {
			t.Errorf("parse(%q) returned nil err, want non-nil", in)
		}
	}
}

func TestParseNegativeRange(t *testing.T) {
	// endByte < startByte should be flagged.
	if _, _, _, _, _, err := parse("[5.0:3.0]"); err == nil {
		t.Error("parse(end<start) returned nil err")
	}
}

// Integration: easybits.Unmarshal pipes through easyio + easyerrors. Verify
// that a structured roundtrip with multiple fields works on a wired-up reader.
func TestUnmarshalIntegrationWithEasyIOReadFullPath(t *testing.T) {
	type packet struct {
		Header uint8  `bits:"[0.0:1.0]"`
		Code   uint16 `bits:"[1.0:3.0]"`
	}
	raw := []byte{0xAB, 0xCD, 0xEF, 0x00}
	r := easyio.NewEasyReader(bytes.NewReader(raw))
	var got packet
	if err := Unmarshal(r, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Header != 0xAB {
		t.Errorf("Header=%x want 0xAB", got.Header)
	}
	if got.Code != 0xCDEF {
		t.Errorf("Code=%x want 0xCDEF", got.Code)
	}
}

func TestUnmarshalErrorPropagatesThroughEasyErrors(t *testing.T) {
	type s struct {
		A uint16 `bits:"[0.0:2.0]"`
	}
	r := easyio.NewEasyReader(strings.NewReader("a")) // 1 byte
	err1 := Unmarshal(r, &s{})

	combined := easyerrors.HandleMultiError(easyerrors.Simple(), nil, err1)
	if combined == nil {
		t.Errorf("expected combined error to surface unmarshal failure")
	}
	if !errors.Is(combined, err1) && combined != err1 {
		// HandleMultiError returns the first non-nil error directly.
		t.Errorf("combined != err1: %v vs %v", combined, err1)
	}
}
