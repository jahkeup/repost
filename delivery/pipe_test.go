package delivery

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPipeToCommand(t *testing.T) {
	in := []byte("this is some data")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cat")

	out, err := pipeToCommand(bytes.NewBuffer(in), cmd)
	assert.NoError(t, err)
	assert.Equal(t, in, out, "expected to echo the input that was piped in")
}

func TestPipeDeliverer(t *testing.T) {
	var out []byte
	in := []byte("this is some data")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cat")

	var piper commandPiper = func(rd io.Reader, cmd *exec.Cmd) ([]byte, error) {
		var err error
		out, err = pipeToCommand(rd, cmd)
		// internals may change for
		var wrappedOut []byte
		if len(out) > 0 {
			wrappedOut = make([]byte, len(out))
			copy(wrappedOut, out)
		}
		return wrappedOut, err
	}

	d := NewCommandDeliverer(cmd)
	d.piper = piper

	err := d.Deliver(bytes.NewBuffer(in))

	assert.NoError(t, err)
	assert.Equal(t, in, out, "expected to echo the input that was piped in")
}
