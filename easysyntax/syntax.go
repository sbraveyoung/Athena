package syntax

import (
	"context"
	"time"
)

func DoLoop(ctx context.Context, f func(), interval time.Duration) {
	go func() {
		f()
		ticker := time.NewTicker(interval)
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case <-ticker.C:
				go f()
			}

		}
		ticker.Stop()
	}()
}
