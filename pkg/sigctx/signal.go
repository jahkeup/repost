package sigctx

import (
	"context"
	"os"
	"os/signal"
)

// subscribe and unsubscribe are vars for testing.
var subscribe = signal.Notify
var unsubcribe = signal.Stop

// WithCancelSignal handles an os.Signal to cancel a context on an os
// signal like SIGINT.
func WithCancelSignal(parent context.Context, cancelSignals ...os.Signal) (context.Context, context.CancelFunc) {
	actualCtx, actualCancel := context.WithCancel(parent)
	c := make(chan os.Signal, 1)

	stopListening := func() {
		unsubcribe(c)
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
		case <-c:
			sigctxCancel()
			return
		}
	}()

	subscribe(c, cancelSignals...)

	return actualCtx, sigctxCancel
}
