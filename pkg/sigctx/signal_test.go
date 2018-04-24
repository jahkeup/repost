package sigctx

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

func reset() {
	subscribe = signal.Notify
	unsubcribe = signal.Stop
}

func TestWithCancelSignal(t *testing.T) {
	maxTestDuration := time.Second * 1
	testCtx, cancel := context.WithTimeout(context.Background(), maxTestDuration)
	defer cancel()

	signal := os.Kill

	signaler := &TestSignaler{}
	signaler.Setup()
	defer signaler.Teardown()

	ctx, _ := WithCancelSignal(testCtx, signal)

	for {
		select {
		case <-testCtx.Done():
			// test exited on its context, not the signal's its a failure.
			t.Fatal("signal context did not finish.")
		case <-ctx.Done():
			// test complete, the context should have exited here, but
			// verify its not a parented event.
			select {
			case <-testCtx.Done():
				t.Fatal("signal context did not finish")
			default:
			}
			return
		default:
			signaler.SendOnce(signal)
		}
	}
}

type TestSignaler struct {
	once    sync.Once
	signal  chan<- os.Signal
	signals []os.Signal
}

func (ts *TestSignaler) IsListeningTo(sig os.Signal) bool {
	for _, signal := range ts.signals {
		if signal == sig {
			return true
		}
	}

	return false
}

func (ts *TestSignaler) Send(sig os.Signal) {
	ts.signal <- sig
}

func (ts *TestSignaler) SendOnce(sig os.Signal) {
	ts.once.Do(func() {
		ts.Send(sig)
	})
}

func (ts *TestSignaler) Subscribe(c chan<- os.Signal, sigs ...os.Signal) {
	ts.signal = c
	ts.signals = sigs
}

func (ts *TestSignaler) Unsubscribe(c chan<- os.Signal) {
	ts.signal = nil
	ts.signals = ts.signals[:0]
}

func (ts *TestSignaler) Setup() {
	subscribe = ts.Subscribe
	unsubcribe = ts.Unsubscribe
}

func (ts *TestSignaler) Teardown() {
	subscribe = signal.Notify
	unsubcribe = signal.Stop
}
