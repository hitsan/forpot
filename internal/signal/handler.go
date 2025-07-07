package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Handler struct {
	sigChan chan os.Signal
	ctx     context.Context
	cancel  context.CancelFunc
}

func New() *Handler {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return &Handler{
		sigChan: sigChan,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (h *Handler) Context() context.Context {
	return h.ctx
}

func (h *Handler) Wait() {
	sig := <-h.sigChan
	fmt.Printf("\nReceived signal: %s\n", sig)
	fmt.Println("Shutting down gracefully...")
	h.cancel()
}

func (h *Handler) Shutdown() {
	h.cancel()
	signal.Stop(h.sigChan)
	close(h.sigChan)
}

