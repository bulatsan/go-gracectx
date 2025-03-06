package gracectx

import (
	"context"
	"os/signal"
	"syscall"
)

func New() (context.Context, context.CancelFunc) {
	return Wrap(context.Background())
}

func Wrap(parent context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
}
