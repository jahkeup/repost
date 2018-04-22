package handler

import (
	"github.com/jahkeup/repost/delivery"
	"github.com/jahkeup/repost/notification"
)

type Vender interface {
	// Vend returns a deliverer suitable for a given notification.
	Vend(notification.DeliveryNotification) delivery.Deliverer
}

// FuncVender vends a deliverer using a function regardless of the
// notification data (which should be handled by a type specific to
// handling this using the Vend interface).
type FuncVender struct {
	newDeliverer func() delivery.Deliverer
}

// NewFuncVender returns a Vender that returns a Deliverer from a
// called function.
func NewFuncVender(fn func() delivery.Deliverer) *FuncVender {
	return &FuncVender{
		newDeliverer: fn,
	}
}

func (r *FuncVender) Vend(_ notification.DeliveryNotification) delivery.Deliverer {
	if r.newDeliverer != nil {
		return r.newDeliverer()
	}
	panic("deliverer func must be provided")
}
