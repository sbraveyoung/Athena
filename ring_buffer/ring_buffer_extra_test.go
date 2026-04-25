package ring_buffer

import (
	"sync"
	"testing"
	"time"
)

func TestNewRingBufferZeroSize(t *testing.T) {
	if got := NewRingBuffer(0); got != nil {
		t.Errorf("NewRingBuffer(0) = %v, want nil", got)
	}
	if got := NewRingBuffer(-1); got != nil {
		t.Errorf("NewRingBuffer(-1) = %v, want nil", got)
	}
}

func TestBuilderArrayDefault(t *testing.T) {
	c := NewRingBuffer(2).Build()
	c.Append("a")
	c.Append("b")
	if got := c.Get(); got != "a" {
		t.Errorf("Get()=%v want a", got)
	}
}

func TestBuilderListExplicit(t *testing.T) {
	c := NewRingBuffer(2).List().Build()
	c.Append("a")
	c.Append("b")
	if got := c.Get(); got != "a" {
		t.Errorf("Get()=%v want a", got)
	}
}

// Both Array and List variants should evict the oldest entry once size+1 items
// have been appended.
func TestNonBlockingOverwriteOldest(t *testing.T) {
	for _, kind := range []string{TYPE_ARRAY, TYPE_LIST} {
		t.Run(kind, func(t *testing.T) {
			c := NewRingBuffer(3).EvictType(kind).Build()
			for i := 1; i <= 5; i++ {
				c.Append(i)
			}
			// After 5 appends with size=3, the oldest two (1,2) are evicted.
			// Reads should yield 3, 4, 5, then nil.
			want := []interface{}{3, 4, 5, nil}
			for i, w := range want {
				if got := c.Get(); got != w {
					t.Errorf("Get[%d]=%v, want %v", i, got, w)
				}
			}
		})
	}
}

func TestNonBlockingEmptyReturnsNil(t *testing.T) {
	c := NewRingBuffer(3).Build()
	if got := c.Get(); got != nil {
		t.Errorf("Get() on empty = %v, want nil", got)
	}
}

func TestArrayWrapAroundAfterRead(t *testing.T) {
	c := NewRingBuffer(2).Array().Build()
	c.Append(1)
	c.Append(2)
	if got := c.Get(); got != 1 {
		t.Fatalf("Get()=%v want 1", got)
	}
	c.Append(3)
	if got := c.Get(); got != 2 {
		t.Errorf("Get()=%v want 2", got)
	}
	if got := c.Get(); got != 3 {
		t.Errorf("Get()=%v want 3", got)
	}
}

// Block mode is a buffered channel; Get must block until something is appended
// and Append must block (well, drop the oldest) once the channel is full.
func TestBlockingGetWaits(t *testing.T) {
	c := NewRingBuffer(1).Block().Build()
	got := make(chan interface{}, 1)
	go func() { got <- c.Get() }()

	select {
	case <-got:
		t.Fatal("Get returned before any value was appended")
	case <-time.After(20 * time.Millisecond):
	}

	c.Append("hello")
	select {
	case v := <-got:
		if v != "hello" {
			t.Errorf("got %v, want hello", v)
		}
	case <-time.After(time.Second):
		t.Fatal("Get did not unblock after Append")
	}
}

// When the blocking buffer is full, Append should drop the oldest entry instead
// of blocking. Verify by appending size+1 items, then reading.
func TestBlockingAppendDropsOldestWhenFull(t *testing.T) {
	c := NewRingBuffer(2).Block().Build()
	c.Append(1)
	c.Append(2)
	c.Append(3) // capacity exceeded; oldest (1) evicted
	if got := c.Get(); got != 2 {
		t.Errorf("Get()=%v want 2 (after eviction of 1)", got)
	}
	if got := c.Get(); got != 3 {
		t.Errorf("Get()=%v want 3", got)
	}
}

// Block mode is a buffered channel; with a buffer at least as large as the
// total number of items, no eviction occurs and we should observe full
// ordered delivery from a producer to a consumer goroutine.
func TestBlockingProducerConsumer(t *testing.T) {
	const N = 200
	c := NewRingBuffer(N).Block().Build()

	got := make([]int, 0, N)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			got = append(got, c.Get().(int))
		}
	}()

	for i := 1; i <= N; i++ {
		c.Append(i)
	}
	wg.Wait()

	if len(got) != N {
		t.Fatalf("consumer received %d, want %d", len(got), N)
	}
	for i, v := range got {
		if v != i+1 {
			t.Errorf("got[%d]=%d, want %d (order broken)", i, v, i+1)
			break
		}
	}
}
