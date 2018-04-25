package notification

import (
	"context"
	"encoding/json"
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
		if msg.Body == nil || aws.StringValue(msg.Body) == "" {
			err := errors.Errorf("MessageId %q: body was empty", aws.StringValue(msg.MessageId))
			s.log.Error(err)
			continue
		}
		body := []byte(aws.StringValue(msg.Body))
		dn, err := unwrapDelivery(body)
		if err != nil {
			s.log.Errorf("error unwrapping delivery notification: %s", err)
			s.log.Debugf("sqs message: %q", msg.Body)
			continue
		}
		err = s.handler.HandleDelivery(*dn)
		if err != nil {
			s.log.Error(err)
			continue
		}
		s.remove(msg)
	}
	return nil
}

func (s *sqsNotification) remove(msg *sqs.Message) error {
	if msg == nil {
		return errors.New("nil sqs message provided")
	}
	return nil
}

func unwrapDelivery(msgBody []byte) (*DeliveryNotification, error) {
	// sqsBodyEnvelope is the json body of the SNS notification in the sqs
	// body field. Yay double unmarshaling.
	type sqsBodyEnvelope struct {
		Type      string
		MessageId string
		TopicArn  string
		Subject   string
		Message   string
	}
	env := &sqsBodyEnvelope{}
	err := json.Unmarshal(msgBody, env)
	if err != nil {
		return nil, errors.Wrap(err, "unexpected envelope")
	}
	dn := &DeliveryNotification{}
	err = json.Unmarshal([]byte(env.Message), dn)
	if err != nil {
		return nil, errors.Wrap(err, "envelope message wasn't a delivery notification")
	}
	return dn, nil
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
