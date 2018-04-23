package main

import (
	"context"
	"os"

	"github.com/jahkeup/repost/pkg/sigctx"
)

// Context provides the main context that handles signals and such.
func Context() (context.Context, context.CancelFunc) {
	ctx, cancel := sigctx.WithCancelSignal(context.Background(), os.Interrupt, os.Kill)
	return ctx, cancel
}
