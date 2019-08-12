// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package govnr

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

type Errorer interface {
	Error(err error)
}

type ContextEndedChan chan struct{}

// Runs f() in a new goroutine; if it panics, logs the error and stack trace to the specified Errorer
//
// Deprecated; use Once instead
func GoOnce(errorHandler Errorer, f func()) {
	go func() {
		tryOnce(errorHandler, f)
	}()
}

// Runs f() in a new goroutine; if it panics, logs the error and stack trace to the specified Errorer
func Once(errorHandler Errorer, f func()) {
	go func() {
		tryOnce(errorHandler, f)
	}()
}

// Runs f() in a new goroutine; if it panics, logs the error and stack trace to the specified Errorer
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

// Runs f() in a new goroutine; if it panics, logs the error and stack trace to the specified Errorer
// If the provided Context isn't closed, re-runs f()
// Returns a construct allowing to wait for graceful shutdown
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
