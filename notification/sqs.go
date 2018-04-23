package notification

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"golang.org/x/sync/semaphore"
)

const (
	// this is the time between waits for SQS, backoffs could also be
	// used to increase this as demand requires to a ceiling.
	sqsWaitDuration    = 60 * time.Second
	sqsMaxPollDuration = sqsWaitDuration * 2
	sqsMaxInFlight     = 3
)

var (
	ErrNoHandler = errors.New("No handler configured for deliveries")
)

type sqsNotification struct {
	queue    string
	sqs      sqs.SQS
	inflight *semaphore.Weighted
	log      logrus.FieldLogger

	handler DeliveryHandler
}

func (s *sqsNotification) SetDeliveryHandler(h DeliveryHandler) {
	s.handler = h
}

func (s *sqsNotification) Poll(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			pollctx, _ := context.WithTimeout(ctx, sqsMaxPollDuration)
			err := s.poll(pollctx)
			if err != nil {
				s.log.Errorf("SQS Poll returned an error: %s", err)
			}
		}
	}
}

func (s *sqsNotification) poll(ctx context.Context) error {
	if s.handler == nil {
		return ErrNoHandler
	}

	if !s.inflight.TryAcquire(1) {
		s.log.Warnf("Would exceed max inflight message requests (at max: %d), blocking new requests", sqsMaxInFlight)
		acqErr := s.inflight.Acquire(ctx, 1)
		if acqErr != nil {
			// context is finished and acquire failed, bail out.
			return acqErr
		}
		s.log.Info("Resuming handling of messages")
	}

	recv, err := s.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.queue),
		MaxNumberOfMessages:   aws.Int64(10), // max
		WaitTimeSeconds:       aws.Int64(int64(sqsWaitDuration.Seconds())),
		MessageAttributeNames: aws.StringSlice([]string{sqs.QueueAttributeNameAll}),
	})
	if err != nil {
		s.inflight.Release(1)
		return err
	}

	err = s.deliver(recv.Messages)
	s.inflight.Release(1)
	return err
}

func (s *sqsNotification) deliver(msgs []*sqs.Message) error {
	s.log.Debugf("Messages: %#v", msgs)
	return nil
}

func NewSQS(ctx context.Context, client sqs.SQS, queueUrl string) *sqsNotification {
	sema := semaphore.NewWeighted(sqsMaxInFlight)
	notif := sqsNotification{
		queue:    queueUrl,
		inflight: sema,
		log:      logrus.WithField("poller", "sqs"),
	}

	return &notif
}
