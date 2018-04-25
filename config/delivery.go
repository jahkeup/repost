package config

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jahkeup/repost/handler"
	noti "github.com/jahkeup/repost/notification"
)

// Delivery configuration
type Delivery struct {
	// Pipe delivery command
	Pipe PipeDelivery
	// This causes the message to be retained in S3. This may be useful
	// for archival purposes or if you're using IA with lifecycles to
	// age out messages. If retaining an archival copy of messages
	// stored is desirable, inbound rules for SES should be tuned to
	// accommodate such goals - this isn't the tool for that.
	KeepMessages bool
}

// Vender returns the configured Vender.
func (d *Delivery) vender() (handler.Vender, error) {
	return &d.Pipe, nil
}

// Handler is the configured handler for receiving messages.
func (d *Delivery) Handler(ctx context.Context, sess *session.Session) (noti.DeliveryHandler, error) {
	s3 := s3.New(sess)
	configuredVender, err := d.vender()

	if err != nil {
		return nil, err
	}
	h := handler.NewS3(s3, configuredVender).KeepMessages(d.KeepMessages)
	return h, nil
}
