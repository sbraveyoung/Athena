package pool

import (
	"io"
	"sync"
	"testing"
)

// Compile-time assertion: Buffer must implement io.ByteWriter — the WriteByte
// signature regression we just fixed.
var _ io.ByteWriter = (*Buffer)(nil)

func TestBufferWrite(t *testing.T) {
	var b Buffer
	b.Write([]byte("hello"))
	if got := b.String(); got != "hello" {
		t.Errorf("String()=%q want %q", got, "hello")
	}
	b.Write([]byte(" world"))
	if got := b.String(); got != "hello world" {
		t.Errorf("String()=%q want %q", got, "hello world")
	}
	if got := b.Len(); got != len("hello world") {
		t.Errorf("Len()=%d want %d", got, len("hello world"))
	}
}

func TestBufferWriteString(t *testing.T) {
	var b Buffer
	b.WriteString("foo")
	b.WriteString("bar")
	if got := b.String(); got != "foobar" {
		t.Errorf("WriteString chain = %q want foobar", got)
	}
}

func TestBufferWriteByte(t *testing.T) {
	var b Buffer
	for _, c := range []byte("abc") {
		if err := b.WriteByte(c); err != nil {
			t.Fatalf("WriteByte(%q) err: %v", c, err)
		}
	}
	if got := b.String(); got != "abc" {
		t.Errorf("WriteByte sequence = %q want abc", got)
	}
}

func TestBufferUsedAsByteWriter(t *testing.T) {
	// Use Buffer through the io.ByteWriter contract; this exercises the fact
	// that the standard library will silently downcast to that interface.
	var b Buffer
	var w io.ByteWriter = &b
	for _, c := range []byte("xyz") {
		if err := w.WriteByte(c); err != nil {
			t.Fatalf("ByteWriter.WriteByte: %v", err)
		}
	}
	if got := b.String(); got != "xyz" {
		t.Errorf("got %q want xyz", got)
	}
}

func TestGetBufferReturnsEmpty(t *testing.T) {
	buf, put := GetBuffer()
	defer put(buf)
	if buf.Len() != 0 {
		t.Errorf("freshly Get'd buffer len=%d, want 0", buf.Len())
	}
}

func TestGetBufferReusesUnderlyingArray(t *testing.T) {
	// Put a buffer back, Get one, and verify we are likely getting the same
	// underlying array (capacity preserved). sync.Pool is allowed to drop
	// items so we retry a few times — flakiness here would manifest as a
	// non-failure, not a false positive.
	first, put := GetBuffer()
	first = append(first, []byte("seed-data-to-grow-cap")...)
	cap1 := cap(first)
	put(first)

	for i := 0; i < 50; i++ {
		got, putGot := GetBuffer()
		if cap(got) == cap1 {
			putGot(got)
			return
		}
		putGot(got)
	}
	t.Skip("sync.Pool dropped the buffer; cannot prove reuse on this run")
}

func TestGetBufferIsCleared(t *testing.T) {
	buf, put := GetBuffer()
	buf = append(buf, []byte("dirty")...)
	put(buf)

	buf2, put2 := GetBuffer()
	defer put2(buf2)
	if buf2.Len() != 0 {
		t.Errorf("re-Get'd buffer not cleared, len=%d", buf2.Len())
	}
}

func TestGetSlice(t *testing.T) {
	s, put := GetSlice()
	defer put(s)
	if len(s) != 0 {
		t.Errorf("len(slice)=%d want 0", len(s))
	}
	if cap(s) <= 0 {
		t.Errorf("expected non-zero cap, got %d", cap(s))
	}
}

func TestJoinEmpty(t *testing.T) {
	out := Join(nil, []byte(","))
	if len(out) != 0 {
		t.Errorf("Join(nil) = %q want empty", out)
	}
}

func TestJoinSingle(t *testing.T) {
	in := []Buffer{Buffer("only")}
	out := Join(in, []byte(","))
	if string(out) != "only" {
		t.Errorf("Join(single)=%q want %q", out, "only")
	}
	// Verify the result is a copy: mutating the input must not affect the result.
	in[0][0] = 'X'
	if string(out) != "only" {
		t.Errorf("Join(single) returned aliased slice; got %q after mutation", out)
	}
}

func TestJoinMulti(t *testing.T) {
	in := []Buffer{Buffer("a"), Buffer("bb"), Buffer("ccc")}
	out := Join(in, []byte(", "))
	if string(out) != "a, bb, ccc" {
		t.Errorf("Join multi = %q", out)
	}
}

func TestJoinEmptySeparator(t *testing.T) {
	in := []Buffer{Buffer("foo"), Buffer("bar")}
	out := Join(in, nil)
	if string(out) != "foobar" {
		t.Errorf("Join with nil sep = %q", out)
	}
}

// The pools are concurrent-safe; this test stresses Get/Put under contention so
// `go test -race` can catch any future regression.
func TestPoolConcurrent(t *testing.T) {
	const goroutines = 16
	const ops = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				buf, put := GetBuffer()
				buf.WriteString("x")
				put(buf)

				slice, putS := GetSlice()
				slice = append(slice, Buffer("y"))
				putS(slice)
			}
		}()
	}
	wg.Wait()
}
