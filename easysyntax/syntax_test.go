package easysyntax

import (
	"context"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

// Regression for the busy-loop bug where break inside select did not exit the
// outer for loop. After ctx is cancelled the goroutine must terminate, so
// runtime.NumGoroutine should return to its baseline.
func TestDoLoopExitsOnContextCancel(t *testing.T) {
	before := runtime.NumGoroutine()

	ctx, cancel := context.WithCancel(context.Background())
	var calls int32
	DoLoop(ctx, func() { atomic.AddInt32(&calls, 1) }, time.Millisecond)

	// Let the loop tick at least once.
	time.Sleep(10 * time.Millisecond)
	cancel()

	// Give the goroutine a moment to observe the cancellation and exit.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= before {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if got := runtime.NumGoroutine(); got > before {
		t.Errorf("goroutine leaked: before=%d after=%d", before, got)
	}

	// And verify no further f() calls happen after cancel.
	stoppedAt := atomic.LoadInt32(&calls)
	time.Sleep(20 * time.Millisecond)
	if got := atomic.LoadInt32(&calls); got != stoppedAt {
		t.Errorf("f() kept being called after cancel: stoppedAt=%d after=%d", stoppedAt, got)
	}
}

func TestDoLoopRunsImmediatelyAndPeriodically(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var calls int32
	DoLoop(ctx, func() { atomic.AddInt32(&calls, 1) }, 10*time.Millisecond)

	// Synchronous first call from DoLoop happens before the goroutine even
	// starts, so by the time DoLoop returns we should already see >= 1.
	if got := atomic.LoadInt32(&calls); got < 1 {
		t.Fatalf("expected >=1 immediate call, got %d", got)
	}

	time.Sleep(55 * time.Millisecond)
	if got := atomic.LoadInt32(&calls); got < 3 {
		t.Errorf("expected several ticks within 55ms, got %d", got)
	}
}
