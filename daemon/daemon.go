package daemon

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/jahkeup/repost/config"
	noti "github.com/jahkeup/repost/notification"
	"github.com/pkg/errors"
)

type Daemon struct {
	log    logrus.FieldLogger
	poller noti.Poller
}

// New configures a daemon for running as configured.
func New(ctx context.Context, config *config.Config) (*Daemon, error) {
	config.General.Apply()

	log := logrus.WithField("context", "daemon")
	session, err := config.Session()
	if err != nil {
		err = errors.Wrap(err, "could not get configured session")
		log.Debug(err)
		return nil, err
	}
	handler, err := config.Delivery.Handler(ctx, session)
	if err != nil {
		err = errors.Wrap(err, "could not get configured handler")
		log.Debug(err)
		return nil, err
	}
	poller, err := config.Notification.Poller(ctx, session)
	if err != nil {
		err = errors.Wrap(err, "could not get configured poller")
		log.Debug(err)
		return nil, err
	}

	if cpoller, ok := poller.(noti.ConfigurablePoller); ok {
		cpoller.SetDeliveryHandler(handler)
	} else {
		log.Warn("Configured poller does not allow configurable handler")
	}

	d := &Daemon{poller: poller, log: log}

	err = d.logCallerIdentity(session)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Run blocks and runs the poller.
func (d *Daemon) Run(ctx context.Context) error {
	d.log.Info("running")
	for {
		select {
		case <-ctx.Done():
			d.log.Debug("context has finished")
			d.log.Info("shutting down")
			return nil
		default:
			d.log.Debug("polling")
			err := d.poller.Poll(ctx)
			d.log.Debug("poll complete")
			if err != nil {
				d.log.Debug("poller errored, exiting run loop")
				return err
			}
		}
	}
}
