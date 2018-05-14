package config

import (
	"bytes"
	"context"
	"os/exec"
	"text/template"
	"time"

	"github.com/jahkeup/repost/delivery"
	noti "github.com/jahkeup/repost/notification"
	"github.com/pkg/errors"
)

const (
	PipeDeliveryDefaultTimeout = time.Second * 10
)

// PipeDelivery configures a handler to deliver via a command and
// stdin.
type PipeDelivery struct {
	// Command is a go template that can use details in the message
	// delivery notification.
	Command string
	Timeout time.Duration

	t *template.Template
}

func (p *PipeDelivery) template() (*template.Template, error) {
	if p.t != nil {
		return p.t, nil
	}

	t := template.New("pipe-delivery").Funcs(funcMap)
	parsed, err := t.Parse(p.Command)
	if err != nil {
		return nil, err
	}
	p.t = parsed
	return p.t, nil
}

func (p *PipeDelivery) cmd(dn noti.DeliveryNotification) (*exec.Cmd, error) {
	if p.Timeout == 0 {
		p.Timeout = PipeDeliveryDefaultTimeout
	}

	t, err := p.template()
	if err != nil {
		return nil, errors.Wrap(err, "error reading template string")
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, dn)
	if err != nil {
		return nil, errors.Wrap(err, "error executing template string")
	}

	ctx, _ := context.WithTimeout(context.Background(), p.Timeout)
	cmd := exec.CommandContext(ctx, "sh", "-c", buf.String())
	return cmd, nil
}

func (p *PipeDelivery) Vend(dn noti.DeliveryNotification) (delivery.Deliverer, error) {
	// Command is a template that can have some details placed into it
	// that are present in the original notification. This is a really
	// small subset. So much so that it might be worth delaying this
	// handling and buffering locally.
	cmd, err := p.cmd(dn)
	if err != nil {
		return nil, err
	}
	deliverer := delivery.NewCommandDeliverer(cmd)
	return deliverer, nil
}
