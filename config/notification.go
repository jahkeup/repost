package config

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	noti "github.com/jahkeup/repost/notification"
	"github.com/pkg/errors"
)

// Notification configuration
type Notification struct {
	// QueueURL is the url of the SQS Queue to receive messages from.
	QueueURL string
}

// Poller is limited to the context and should be polled by the
// daemon.
func (n *Notification) Poller(_ context.Context, sess *session.Session) (noti.Poller, error) {
	sqs := sqs.New(sess)
	if sqs == nil {
		return nil, errors.New("error creating sqs client from config")
	}

	return noti.NewSQS(sqs, n.QueueURL), nil
}
