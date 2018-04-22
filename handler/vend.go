package handler

import (
	"github.com/jahkeup/repost/delivery"
	"github.com/jahkeup/repost/notification"
)

type Vender interface {
	// Vend returns a deliverer suitable for a given notification.
	Vend(notification.DeliveryNotification) delivery.Deliverer
}

type ReusableVender struct {
	d delivery.Deliverer
}

func (r *ReusableVender) Vend(_ notification.DeliveryNotification) delivery.Deliverer {
	return r.d
}
