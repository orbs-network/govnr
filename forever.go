package govnr

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

// Runs f() in a new goroutine; if it panics, emits the error to the specified Errorer
// If the provided Context isn't closed, re-runs f()
// Returns a channel that is closed when the goroutine has quit due to context ending
//
// Deprecated; use Forever instead
func GoForever(ctx context.Context, errorHandler Errorer, f func()) ContextEndedChan {
	c := make(ContextEndedChan)
	go func() {
		defer close(c)

		for {
			tryOnce(errorHandler, f)
			//TODO(v1) report number of restarts to metrics
			if ctx.Err() != nil { // this returns non-nil when context has been closed via cancellation or timeout or whatever
				return
			}
		}
	}()
	return c
}

type ForeverHandle struct {
	sync.Mutex
	closed       chan struct{}
	errorHandler Errorer
	name         string
	supervised   bool
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
	h.Lock()
	defer h.Unlock()
	h.supervised = true
}

func (h *ForeverHandle) terminated() {
	close(h.closed)
	h.Lock()
	defer h.Unlock()
	if !h.supervised {
		h.errorHandler.Error(errors.Errorf("Forever governed goroutine %s terminated without being supervised", h.name))
	}
}

// Runs f() in a new goroutine; if it panics, emits the error to the provided Errorer.
// If the provided Context isn't closed, re-runs f().
// Returns a ForeverHandle to allow a Supervisor to wait for graceful shutdown.
// When f() exists normally, if the ForeverHandle hasn't been passed to a Supervisor, an error will be emitted to the provided Errorer.
func Forever(ctx context.Context, name string, errorHandler Errorer, f func()) *ForeverHandle {
	h := &ForeverHandle{closed: make(ContextEndedChan), name: name, errorHandler: errorHandler}
	go func() {
		defer h.terminated()

		for {
			tryOnce(errorHandler, f)
			if ctx.Err() != nil { // this returns non-nil when context has been closed via cancellation or timeout or whatever
				return
			}
		}
	}()
	return h
}

