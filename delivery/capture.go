package delivery

import (
	"bytes"
	"io"

	"github.com/Sirupsen/logrus"
)

type Capture struct {
	messageData bytes.Buffer
}

func NewCapture() *Capture {
	return &Capture{}
}

func (c *Capture) Data() []byte {
	return c.messageData.Bytes()
}

func (c *Capture) Deliver(rd io.Reader) error {
	if c.messageData.Len() > 0 {
		return ErrDelivererAlreadyUsed
	}

	log := logrus.WithField("deliverer", "Capture")

	n, err := io.Copy(&c.messageData, rd)
	if err != nil {
		return err
	}
	closed := closeIfCloser(rd)
	if closed {
		log.Debug("closing delivery reader")
	}

	log.Debugf("Read %d bytes from delivery reader", n)
	log.Debugf("Message data: %q", c.messageData.Bytes())

	return nil
}
