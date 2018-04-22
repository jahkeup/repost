package handler

import (
	noti "github.com/jahkeup/repost/notification"
)

type Mock struct {
	OnHandleDelivery func(noti.DeliveryNotification) error
}

func (n *Mock) HandleDelivery(dn noti.DeliveryNotification) error {
	if n.OnHandleDelivery != nil {
		return n.OnHandleDelivery(dn)
	}

	return nil
}
