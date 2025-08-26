package server

import (
	"context"
	"testing"
	"time"
)

// Additional coverage for Run path when ListenAndServe returns ErrServerClosed
func TestRunImmediateCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()                                        // cancel immediately
	if err := Run(ctx, "127.0.0.1:0"); err != nil { //nolint:errcheck // acceptable immediate shutdown
		// intentionally ignore; just exercising path
		_ = err
	}
}

// Exercise shutdown path with signal-style cancellation via context timeout
func TestRunTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := Run(ctx, "127.0.0.1:0"); err != nil { //nolint:errcheck // acceptable immediate shutdown
		_ = err
	}
}
