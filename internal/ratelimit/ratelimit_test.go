package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/driftwatch/internal/ratelimit"
)

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Config{RequestsPerSecond: 0})
	if err == nil {
		t.Fatal("expected error for RequestsPerSecond=0, got nil")
	}
}

func TestNew_ValidRate(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{RequestsPerSecond: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()
}

func TestWait_AcquiresToken(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{RequestsPerSecond: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()

	ctx := context.Background()
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("Wait returned unexpected error: %v", err)
	}
}

func TestWait_CancelledContext(t *testing.T) {
	// Rate of 1 rps; pre-fill is consumed immediately so next Wait should block.
	l, err := ratelimit.New(ratelimit.Config{RequestsPerSecond: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()

	ctx := context.Background()
	// Drain the pre-filled token.
	if err := l.Wait(ctx); err != nil {
		t.Fatalf("first Wait failed: %v", err)
	}

	// Now the bucket is empty; cancel context immediately.
	cancelCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	err = l.Wait(cancelCtx)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

func TestWait_RefillsOverTime(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Config{RequestsPerSecond: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()

	ctx := context.Background()

	// Drain all pre-filled tokens.
	for i := 0; i < 10; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("Wait %d failed: %v", i, err)
		}
	}

	// Wait for at least one refill tick (rate = 100ms per token at 10 rps).
	time.Sleep(150 * time.Millisecond)

	ctx2, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	if err := l.Wait(ctx2); err != nil {
		t.Fatalf("expected token after refill, got error: %v", err)
	}
}
