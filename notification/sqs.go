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
	sqsWaitDuration    = 20 * time.Second
	sqsMaxPollDuration = sqsWaitDuration * 2
	sqsMaxInFlight     = 3
)

var (
	ErrNoHandler = errors.New("No handler configured for deliveries")
)

// sqsReceiver describes the methods called to receive messages from a
// SQS queue.
type SQSReceiver interface {
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
}

type sqsNotification struct {
	queue    string
	sqs      SQSReceiver
	inflight *semaphore.Weighted
	log      logrus.FieldLogger

	handler DeliveryHandler
}

func (s *sqsNotification) SetDeliveryHandler(h DeliveryHandler) {
	s.handler = h
}

func (s *sqsNotification) Poll(ctx context.Context) error {
	return s.poll(ctx)
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return nil
	// 	default:
	// 		pollctx, _ := context.WithTimeout(ctx, sqsMaxPollDuration)
	// 		err := s.poll(pollctx)
	// 		if err != nil {
	// 			s.log.Errorf("SQS Poll returned an error: %s", err)
	// 		}
	// 	}
	// }
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

	maxMessages := int64(10)
	s.log.Debugf("Polling SQS for new %d message(s)", maxMessages)
	recv, err := s.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.queue),
		MaxNumberOfMessages:   aws.Int64(maxMessages), // max
		WaitTimeSeconds:       aws.Int64(int64(sqsWaitDuration.Seconds())),
		MessageAttributeNames: aws.StringSlice([]string{sqs.QueueAttributeNameAll}),
	})
	if err != nil {
		s.inflight.Release(1)
		return err
	}

	if recv == nil {
		s.inflight.Release(1)
		return errors.New("sqs response was empty")
	}

	err = s.deliver(recv.Messages)
	s.inflight.Release(1)
	return err
}

func (s *sqsNotification) deliver(msgs []*sqs.Message) error {
	for _, msg := range msgs {
		dn, err := sqsMessageToDeliveryNotification(msg)
		if err != nil {
			s.log.Errorf("error extracting delivery notification from sqs message: %s", err)
			if msg.Body != nil {
				s.log.Debugf("sqs message: %q", msg.Body)
			}
			continue
		}
		if dn == nil {
			errors.New("extracted delivery notification was empty")
			s.log.Error(err)
			return err
		}
		s.handler.HandleDelivery(*dn)
	}
	return nil
}

func sqsMessageToDeliveryNotification(msg *sqs.Message) (*DeliveryNotification, error) {
	logrus.Printf("%s", aws.StringValue(msg.Body))
	return nil, nil
}

func NewSQS(receiver SQSReceiver, queueUrl string) *sqsNotification {
	sema := semaphore.NewWeighted(sqsMaxInFlight)
	notif := sqsNotification{
		sqs:      receiver,
		queue:    queueUrl,
		inflight: sema,
		log: logrus.WithFields(logrus.Fields{
			"context": "poller",
			"poller":  "sqs",
		}),
	}

	return &notif
}
