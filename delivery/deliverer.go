package delivery

import (
	"io"
)

// Deliverer can deliver message readers to destinations.
type Deliverer interface {
	// Deliver a message reader.
	Deliver(msgRd io.Reader) error
}

func closeIfCloser(rd io.Reader) bool {
	if closer, ok := rd.(io.ReadCloser); ok {
		_ = closer.Close()
		return true
	}
	return false
}
