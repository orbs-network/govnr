package govnr

import (
	"context"
	"github.com/pkg/errors"
)

type ForeverHandle struct {
	closed       chan struct{}
	supervised   chan struct{}
	errorHandler Errorer
	name         string
}

func (h *ForeverHandle) WaitUntilShutdown(timeoutCtx context.Context) {
	select {
	case <-h.closed:
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			h.errorHandler.Error(errors.Wrapf(timeoutCtx.Err(), "Forever governed goroutine %s timed out while waiting for shutdown", h.name))
		}
	}
}

func (h *ForeverHandle) Done() ContextEndedChan {
	return h.closed
}

func (h *ForeverHandle) MarkSupervised() {
	close(h.supervised)
}

func (h *ForeverHandle) waitUntilSupervised(ctx context.Context) {
	select {
	case <-h.supervised:
	case <-ctx.Done():
	}
}

func (h *ForeverHandle) terminated() {
	close(h.closed)
}

// Runs f() in a new goroutine; if it panics, emits the error to the provided Errorer.
// If the provided Context isn't closed, re-runs f().
// Returns a ForeverHandle to allow a Supervisor to wait for graceful shutdown.
// When f() exists normally, if the ForeverHandle hasn't been passed to a Supervisor, an error will be emitted to the provided Errorer.
func Forever(ctx context.Context, name string, errorHandler Errorer, f func()) *ForeverHandle {
	h := &ForeverHandle{closed: make(ContextEndedChan), supervised: make(chan struct{}), name: name, errorHandler: errorHandler}
	go func() {
		h.waitUntilSupervised(ctx)
		defer h.terminated()
		for ctx.Err() == nil { // this will break when context has been closed via cancellation or timeout or whatever
			tryOnce(errorHandler, f)
		}
	}()
	return h
}
