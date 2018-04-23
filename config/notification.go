package config

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	noti "github.com/jahkeup/repost/notification"
)

// Notification configuration
type Notification struct {
	// QueueURL is the url of the SQS Queue to receive messages from.
	QueueURL string
}

// Poller is limited to the context and should be polled by the
// daemon.
func (n *Notification) Poller(ctx context.Context, sess *session.Session) (noti.Poller, error) {
	sqs := sqs.New(sess)
	return noti.NewSQS(ctx, sqs, n.QueueURL), nil
}
