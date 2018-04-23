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

type ConfigurablePoller interface {
	Poller
	SetDeliveryHandler(h DeliveryHandler)
}
