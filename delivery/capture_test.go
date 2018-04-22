package delivery

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaptureDeliverer(t *testing.T) {
	in := []byte("this is some data")

	capture := NewCapture()

	err := capture.Deliver(bytes.NewBuffer(in))
	require.NoError(t, err)
	out := capture.Data()
	assert.Equal(t, in, out)
}
