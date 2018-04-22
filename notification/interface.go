package notification

import (
	"context"
)

type DeliveryHandler interface {
	HandleDelivery(DeliveryNotification) error
}

type Poller interface {
	Poll(context.Context) error
}
