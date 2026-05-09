package easysyntax

import (
	"context"
	"time"
)

func DoLoop(ctx context.Context, f func(), interval time.Duration) {
	f()
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				f()
			}
		}
	}()
}
