package notified

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ses"
)

// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/receiving-email-notifications-contents.html
type DeliveryNotification struct {
	NotificationType string
	Receipt          ses.ReceiptAction
	Mail             mailObject
}

// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/receiving-email-notifications-contents.html#receiving-email-notifications-contents-mail-object
type mailObject struct {
	Destination      []string
	Source           string
	MessageId        string
	Timestamp        time.Time
	Headers          []mailHeader
	CommonHeaders    []mailHeader
	HeadersTruncated bool
}

type mailHeader struct {
	Name  string
	Value string
}
