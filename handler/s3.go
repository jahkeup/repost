package handler

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	noti "github.com/jahkeup/repost/notification"
	"github.com/pkg/errors"
)

type S3 struct {
	client s3.S3
	vender Vender
}

func NewS3(client s3.S3, vender Vender) *S3 {
	return &S3{
		client: client,
		vender: vender,
	}
}

func (s *S3) HandleDelivery(n noti.DeliveryNotification) error {
	return nil
}

func (s *S3) messageBucketObject(n noti.DeliveryNotification) (string, string, error) {
	if n.NotificationType != "delivery" {
		return "", "", errors.Errorf("cannot handle notification type: %q", n.NotificationType)
	}

	if n.Receipt.S3Action == nil {
		return "", "", errors.Errorf("no S3Action taken in receipt of message, cannot handle")
	}

	s3md := n.Receipt.S3Action
	bucket := aws.StringValue(s3md.BucketName)
	prefix := aws.StringValue(s3md.ObjectKeyPrefix)

	objectKey := fmt.Sprintf("%s/%s", prefix, n.Mail.MessageId)

	return bucket, objectKey, nil
}
