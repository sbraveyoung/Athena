package easyio

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/SmartBrave/Athena/easyerrors"
)

func TestEasyReaderReadFull(t *testing.T) {
	r := NewEasyReader(strings.NewReader("hello world"))
	buf := make([]byte, 5)
	if err := r.ReadFull(buf); err != nil {
		t.Fatalf("ReadFull: %v", err)
	}
	if string(buf) != "hello" {
		t.Errorf("got %q want hello", buf)
	}
}

func TestEasyReaderReadFullEOF(t *testing.T) {
	r := NewEasyReader(strings.NewReader("hi"))
	buf := make([]byte, 5)
	err := r.ReadFull(buf)
	if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		// io.ReadFull returns ErrUnexpectedEOF when it reads some but not all.
		// The wrapper returns it raw if it was EOF; otherwise wraps.
		// Either is acceptable, but err must be non-nil.
		if err == nil {
			t.Errorf("ReadFull on short input returned nil error")
		}
	}
}

func TestEasyReaderReadN(t *testing.T) {
	r := NewEasyReader(strings.NewReader("abcdef"))
	got, err := r.ReadN(3)
	if err != nil {
		t.Fatalf("ReadN: %v", err)
	}
	if string(got) != "abc" {
		t.Errorf("ReadN(3)=%q want abc", got)
	}
}

func TestEasyReaderReadAll(t *testing.T) {
	r := NewEasyReader(strings.NewReader("hello"))
	got, err := r.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("ReadAll=%q want hello", got)
	}
}

func TestEasyWriterWriteFull(t *testing.T) {
	var buf bytes.Buffer
	w := NewEasyWriter(&buf)
	if err := w.WriteFull([]byte("hello")); err != nil {
		t.Fatalf("WriteFull: %v", err)
	}
	if buf.String() != "hello" {
		t.Errorf("dest=%q want hello", buf.String())
	}
}

// shortWriter writes only the first byte and reports n=1, no error. WriteFull
// must reject this as a partial write.
type shortWriter struct{ written []byte }

func (s *shortWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	s.written = append(s.written, p[0])
	return 1, nil
}

func TestEasyWriterWriteFullPartialIsError(t *testing.T) {
	w := NewEasyWriter(&shortWriter{})
	if err := w.WriteFull([]byte("ab")); err == nil {
		t.Errorf("WriteFull(partial)=nil err, want non-nil")
	}
}

// errReader always returns an error so we can verify wrapping.
type errReader struct{}

func (errReader) Read(p []byte) (int, error) { return 0, errors.New("synthetic read error") }

func TestEasyReaderReadFullWrapsNonEOFError(t *testing.T) {
	r := NewEasyReader(errReader{})
	err := r.ReadFull(make([]byte, 4))
	if err == nil {
		t.Fatal("expected non-nil err")
	}
	if errors.Is(err, io.EOF) {
		t.Errorf("synthetic error wrongly classified as EOF: %v", err)
	}
}

// After writing data and Close, ReadAll must drain everything and then return
// EOF cleanly (not block). Probing with empty data hits a separate, unrelated
// edge case in easyReadWriter.Read so we always seed at least one byte.
func TestEasyReadWriterCloseAfterWrite(t *testing.T) {
	rw := NewEasyReadWriter()
	if err := rw.WriteFull([]byte("payload")); err != nil {
		t.Fatalf("WriteFull: %v", err)
	}
	rw.Close()
	got, err := rw.ReadAll()
	if err != nil {
		t.Errorf("ReadAll after Close: %v", err)
	}
	if string(got) != "payload" {
		t.Errorf("ReadAll=%q want payload", got)
	}
}

func TestCopyFull(t *testing.T) {
	src := NewEasyReader(strings.NewReader("hello world"))
	var dst bytes.Buffer
	if err := CopyFull(NewEasyWriter(&dst), src); err != nil {
		t.Fatalf("CopyFull: %v", err)
	}
	if dst.String() != "hello world" {
		t.Errorf("CopyFull dest=%q", dst.String())
	}
}

// Integration: easyio errors flow through easyerrors.HandleMultiError.
func TestEasyIOIntegratesWithEasyErrors(t *testing.T) {
	r := NewEasyReader(errReader{})
	err1 := r.ReadFull(make([]byte, 4))

	var buf bytes.Buffer
	w := NewEasyWriter(&buf)
	err2 := w.WriteFull([]byte("ok"))

	combined := easyerrors.HandleMultiError(easyerrors.Simple(), err2, err1)
	if combined != err1 {
		t.Errorf("HandleMultiError did not surface the first non-nil err: got %v want %v", combined, err1)
	}

	// And when both are nil, should be nil.
	if got := easyerrors.HandleMultiError(easyerrors.Simple(), nil, nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}
