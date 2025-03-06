package gracectx

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx, cancel := New()
	defer cancel()

	if ctx == nil {
		t.Fatal("expected non-nil context")
	}

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done yet")
	default:
	}

	p := syscall.Getpid()

	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(p, syscall.SIGTERM)
	}()

	select {
	case <-ctx.Done():
	case <-time.After(1 * time.Second):
		t.Fatal("context was not canceled after signal")
	}
}

func TestWrap(t *testing.T) {
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	ctx, cancel := Wrap(parentCtx)
	defer cancel()

	if ctx == nil {
		t.Fatal("expected non-nil context")
	}

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done yet")
	default:
	}

	t.Run("cancel_on_signal", func(t *testing.T) {
		localCtx, localCancel := Wrap(context.Background())
		defer localCancel()

		p := syscall.Getpid()

		go func() {
			time.Sleep(100 * time.Millisecond)
			syscall.Kill(p, syscall.SIGTERM)
		}()

		select {
		case <-localCtx.Done():
		case <-time.After(1 * time.Second):
			t.Fatal("context was not canceled after signal")
		}
	})

	t.Run("cancel_on_parent_cancel", func(t *testing.T) {
		parentCtx, parentCancel := context.WithCancel(context.Background())

		wrappedCtx, wrappedCancel := Wrap(parentCtx)
		defer wrappedCancel()

		go func() {
			time.Sleep(100 * time.Millisecond)
			parentCancel()
		}()

		select {
		case <-wrappedCtx.Done():
		case <-time.After(1 * time.Second):
			t.Fatal("wrapped context was not canceled after parent was canceled")
		}
	})
}
