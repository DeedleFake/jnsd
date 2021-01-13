package cli

import (
	"context"
	"os"
	"os/signal"
)

// SignalContext returns a context that is canceled when the given
// signal is received.
func SignalContext(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		c := make(chan os.Signal, 1)
		signal.Notify(c, signals...)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
		}
	}()

	return ctx, cancel
}
