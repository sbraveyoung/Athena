package broadcast

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// startReader runs r.Read() in a goroutine and forwards each (value, alive)
// event onto a channel. The first time Read returns alive=false the channel is
// closed. This matches the canonical user pattern: stop reading on alive=false.
func startReader(r *BroadcastReader) <-chan interface{} {
	ch := make(chan interface{}, 64)
	go func() {
		defer close(ch)
		for {
			p, alive := r.Read()
			if !alive {
				if p != nil {
					ch <- p
				}
				return
			}
			ch <- p
		}
	}()
	return ch
}

// collect drains ch until it closes or the deadline elapses.
func collect(ch <-chan interface{}, deadline time.Duration) []interface{} {
	out := []interface{}{}
	timer := time.NewTimer(deadline)
	defer timer.Stop()
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		case <-timer.C:
			return out
		}
	}
}

func TestWriteIgnoredWhileMetaUnfilled(t *testing.T) {
	bd := NewBroadcast(2)
	bd.Write("ignored-1")
	bd.Write("ignored-2")
	if got := len(bd.data); got != 0 {
		t.Errorf("data len=%d, want 0 (writes before meta full must be dropped)", got)
	}
	if got := bd.wIndex; got != 0 {
		t.Errorf("wIndex=%d, want 0", got)
	}
}

func TestWriteMetaCappedByMetaNum(t *testing.T) {
	bd := NewBroadcast(2)
	bd.WriteMeta("a")
	bd.WriteMeta("b")
	bd.WriteMeta("c") // dropped
	if got := len(bd.meta); got != 2 {
		t.Errorf("meta len=%d, want 2 (overflow should be dropped)", got)
	}
}

func TestMetaThenDataDelivered(t *testing.T) {
	bd := NewBroadcast(2)
	r := NewBroadcastReader(bd)
	ch := startReader(r)

	bd.WriteMeta("m0")
	bd.WriteMeta("m1")
	bd.Write("d0")
	bd.Write("d1")

	got := collect(ch, 200*time.Millisecond)
	// We expect at minimum 2 metas + 2 data values, in that order.
	if len(got) < 4 {
		bd.DisAlive()
		t.Fatalf("got %d events, want >=4: %v", len(got), got)
	}
	if got[0] != "m0" || got[1] != "m1" {
		t.Errorf("metas not delivered first: %v", got[:2])
	}
	if got[2] != "d0" || got[3] != "d1" {
		t.Errorf("data not delivered after meta: %v", got[2:4])
	}
	bd.DisAlive()
}

func TestMultipleReadersGetSameStream(t *testing.T) {
	bd := NewBroadcast(1)
	ch1 := startReader(NewBroadcastReader(bd))
	ch2 := startReader(NewBroadcastReader(bd))

	bd.WriteMeta("m")
	for i := 0; i < 5; i++ {
		bd.Write(i)
	}

	g1 := collect(ch1, 300*time.Millisecond)
	g2 := collect(ch2, 300*time.Millisecond)
	bd.DisAlive()

	for i, got := range [][]interface{}{g1, g2} {
		if len(got) < 6 { // 1 meta + 5 data
			t.Errorf("reader%d got %d events, want >=6: %v", i+1, len(got), got)
		}
		if len(got) > 0 && got[0] != "m" {
			t.Errorf("reader%d first event = %v, want meta 'm'", i+1, got[0])
		}
	}
}

func TestDisAliveWakesBlockedReader(t *testing.T) {
	bd := NewBroadcast(1)
	r := NewBroadcastReader(bd)

	var done int32
	go func() {
		_, _ = r.Read()
		atomic.StoreInt32(&done, 1)
	}()

	time.Sleep(20 * time.Millisecond)
	if atomic.LoadInt32(&done) != 0 {
		t.Fatal("reader returned before any meta was written")
	}

	bd.WriteMeta("m")
	bd.DisAlive()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&done) == 1 {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("reader did not wake after DisAlive")
}

// Regression: previously Read could panic with index 0 when reset was invoked
// but no data had ever been written. Now it must observe alive=false.
func TestReadDoesNotPanicWhenDataEmptyAfterReset(t *testing.T) {
	bd := NewBroadcast(1)
	r := NewBroadcastReader(bd)

	bd.WriteMeta("only-meta")
	bd.Reset()
	bd.DisAlive()

	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Read panicked: %v", rec)
		}
	}()

	ch := startReader(r)
	got := collect(ch, 500*time.Millisecond)
	if len(got) == 0 {
		t.Fatal("reader produced no events at all")
	}
}

func TestProducerMultipleReadersRace(t *testing.T) {
	bd := NewBroadcast(2)
	const readers = 4
	wg := sync.WaitGroup{}
	wg.Add(readers)
	for i := 0; i < readers; i++ {
		r := NewBroadcastReader(bd)
		go func() {
			defer wg.Done()
			for {
				_, alive := r.Read()
				if !alive {
					return
				}
			}
		}()
	}

	bd.WriteMeta(0)
	bd.WriteMeta(1)
	for i := 0; i < 50; i++ {
		bd.Write(i)
	}
	time.Sleep(20 * time.Millisecond)
	bd.DisAlive()

	wg.Wait()
}
