package sigctx

import (
	"context"
	"os"
	"os/signal"
)

// WithCancelSignal handles an os.Signal to cancel a context on an os
// signal like SIGINT.
func WithCancelSignal(parent context.Context, cancelSignals ...os.Signal) (context.Context, context.CancelFunc) {
	actualCtx, actualCancel := context.WithCancel(parent)
	c := make(chan os.Signal, 1)

	stopListening := func() {
		signal.Stop(c)
	}

	sigctxCancel := func() {
		// cancel should also stop the signal listener
		stopListening()
		actualCancel()
	}

	go func() {
		select {
		case <-actualCtx.Done():
			stopListening()
		}
	}()

	signal.Notify(c, os.Interrupt)

	return actualCtx, sigctxCancel
}
