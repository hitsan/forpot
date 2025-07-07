package signal

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSignalHandler(t *testing.T) {
	handler := New()

	// Test that context is not done initially
	select {
	case <-handler.Context().Done():
		t.Error("Context should not be done initially")
	default:
	}

	// Test graceful shutdown
	go func() {
		time.Sleep(100 * time.Millisecond)
		handler.Shutdown()
	}()

	select {
	case <-handler.Context().Done():
		// Expected
	case <-time.After(1 * time.Second):
		t.Error("Context should be done after Shutdown()")
	}
}

func TestSignalHandlerWithSignal(t *testing.T) {
	handler := New()

	// Test signal handling in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		// Send SIGINT to current process
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGINT)
	}()

	// This should return when signal is received
	done := make(chan struct{})
	go func() {
		handler.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Expected
	case <-time.After(1 * time.Second):
		t.Error("Wait() should return after receiving signal")
	}

	// Context should be done after signal
	select {
	case <-handler.Context().Done():
		if handler.Context().Err() != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", handler.Context().Err())
		}
	default:
		t.Error("Context should be done after signal")
	}
}

