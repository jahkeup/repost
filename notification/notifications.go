package notification

import (
	"time"
)

// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/receiving-email-notifications-contents.html
type DeliveryNotification struct {
	NotificationType string
	Receipt          receiptObject
	Mail             mailObject
}

// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/receiving-email-notifications-contents.html#receiving-email-notifications-contents-mail-object
type mailObject struct {
	Destination      []string
	Source           string
	MessageId        string
	Timestamp        time.Time
	Headers          []mailHeader
	CommonHeaders    map[string]interface{}
	HeadersTruncated bool
}

type receiptObject struct {
	Action actionObject
}

type actionObject struct {
	Type            string
	BucketName      string
	ObjectKey       string
	ObjectKeyPrefix string
	KmsKeyArn       string
}

type mailHeader struct {
	Name  string
	Value string
}
