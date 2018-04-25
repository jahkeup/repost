package handler

import (
	"io"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	noti "github.com/jahkeup/repost/notification"
	"github.com/pkg/errors"
)

type s3client interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

type S3 struct {
	client       s3client
	vender       Vender
	keepMessages bool
}

func NewS3(client *s3.S3, vender Vender) *S3 {
	return &S3{
		client: client,
		vender: vender,
	}
}

func (s *S3) KeepMessages(keep bool) *S3 {
	if keep {
		logrus.Info("Handler will NOT remove messages after successful delivery")
	} else {
		logrus.Warn("Handler will remove messages after successful delivery")
	}

	s.keepMessages = keep
	return s
}

func (s *S3) HandleDelivery(n noti.DeliveryNotification) error {
	log := logrus.WithFields(logrus.Fields{
		"message-id": n.Mail.MessageId,
	})

	bucket, object, err := deliveryBucketObject(n)
	if err != nil {
		err = errors.Wrap(err, "couldn't resolve s3 delivery location")
		log.Error(err)
		return err
	}

	out, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	})

	log.Debugf("Fetching message from S3: s3://%s%s", bucket, object)
	if err != nil {
		err = errors.Wrap(err, "could not fetch message from S3")
		log.Error(err)
		return err
	}
	defer drainReadCloser(out.Body)

	// Vend a Deliverer for this notification.
	deliverer, err := s.vender.Vend(n)
	if err != nil {
		return errors.Wrap(err, "deliever unavailable")
	}
	log.Debug("delivering message")
	err = deliverer.Deliver(out.Body)
	if err != nil {
		err = errors.Wrap(err, "could not deliver message")
		log.Error(err)
		return err
	}
	drainReadCloser(out.Body)

	s.remove(bucket, object)

	return nil
}

func (s *S3) remove(bucket string, object string) error {
	if s.keepMessages {
		logrus.Infof("Retaining S3 message: s3://%s/%s", bucket, object)
		return nil
	}
	// TODO: Add S3 object removal
	logrus.Warn("No S3 removal implementation present")
	return nil
}

func drainReadCloser(rc io.ReadCloser) error {
	io.Copy(ioutil.Discard, rc)
	return rc.Close()
}

func deliveryBucketObject(n noti.DeliveryNotification) (string, string, error) {
	if n.NotificationType != "Received" {
		return "", "", errors.Errorf("cannot handle notification type: %q", n.NotificationType)
	}

	if n.Receipt.Action.Type != "S3" {
		return "", "", errors.Errorf("Action taken in receipt of message was not S3, cannot handle")
	}

	s3action := n.Receipt.Action
	switch {
	case s3action.BucketName == "":
		return "", "", errors.New("BucketName was not provided in the action")
	case s3action.ObjectKey == "":
		return "", "", errors.New("ObjectKey was not provided in action")
	case s3action.KmsKeyArn != "":
		return "", "", errors.New("Cannot handle encrypted object")
	}
	return s3action.BucketName, s3action.ObjectKey, nil
}
