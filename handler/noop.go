package handler

import (
	noti "github.com/jahkeup/repost/notification"
)

type Mock struct {
	OnHandleDelivery func(noti.DeliveryNotification)
}

func (n *Mock) HandleDelivery(dn noti.DeliveryNotification) error {

	if n.OnHandleDelivery != nil {
		n.OnHandleDelivery(dn)
	}

	return nil
}
