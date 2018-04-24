package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDaemonHappyPath(t *testing.T) {
	t.Parallel()

	testDuration := 1 * time.Second
	testCtx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Poller should be polled by the daemon at least once while
	// running and should exit without error on its own.
	poller := NewMockPoller(ctrl)
	poller.EXPECT().Poll(gomock.Any()).MinTimes(1)

	daemon := Daemon{
		poller: poller,
		log: logrus.WithFields(logrus.Fields{
			"test": t.Name(),
		}),
	}

	err := daemon.Run(testCtx)
	assert.NoError(t, err)
}
