package config

import (
	"github.com/jahkeup/repost/delivery"
	noti "github.com/jahkeup/repost/notification"
)

// PipeDelivery configures a handler to deliver via a command and
// stdin.
type PipeDelivery struct {
	// Command is a go template that can use details in the message
	// delivery notification.
	Command string
}

func (p *PipeDelivery) Vend(n noti.DeliveryNotification) delivery.Deliverer {
	// Command is a template that can have some details placed into it
	// that are present in the original notification. This is a really
	// small subset. So much so that it might be worth delaying this
	// handling and buffering locally.
	return nil
}
