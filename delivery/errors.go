package delivery

import "github.com/pkg/errors"

var (
	ErrDelivererAlreadyUsed = errors.New("Deliverer cannot be reused")
)
