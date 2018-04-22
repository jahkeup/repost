package delivery

import (
	"bytes"
	"io"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
)

// assert Deliverer interface adherence.
var _ Deliverer = (*commandDeliverer)(nil)

var (
	ErrCommandAlreadyUsed = errors.New("command has already been executed and cannot be reused")
)

// commandPiper is used to qualify the handler executing the commands
// for piping stdin.
type commandPiper func(io.Reader, *exec.Cmd) ([]byte, error)

// CommandDeliverer delivers messages to a command.
type commandDeliverer struct {
	// mu protects the reuse of the command, which _cannot_ happen.
	mu *sync.Mutex
	// cmd is the command that will be executed and delivered to via stdin.
	cmd *exec.Cmd

	// piper is a field here so that tests can override this.
	piper commandPiper
}

// NewCommandDeliverer creates a new Deliverer for delivering a
// message.
func NewCommandDeliverer(cmd *exec.Cmd) *commandDeliverer {
	return &commandDeliverer{
		mu:    &sync.Mutex{},
		cmd:   cmd,
		piper: pipeToCommand,
	}
}

// Deliver message reader to this deliverer.
func (c *commandDeliverer) Deliver(rd io.Reader) error {
	c.mu.Lock()
	if c.cmd == nil {
		c.mu.Unlock()
		return ErrCommandAlreadyUsed
	}
	_, err := c.piper(rd, c.cmd)
	closeIfCloser(rd)
	c.cmd = nil // should not be reused.
	c.mu.Unlock()
	return err
}

// pipeToCommand delivers the input from the reader to a command's stdin.
func pipeToCommand(rd io.Reader, cmd *exec.Cmd) ([]byte, error) {
	outbuf := bytes.NewBuffer(nil)

	cmd.Stdin = rd
	cmd.Stdout = outbuf

	err := cmd.Start()
	if err != nil {
		errors.Wrap(err, "could not start the provided command")
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return outbuf.Bytes(), nil
}
