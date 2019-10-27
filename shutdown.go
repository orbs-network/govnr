package govnr

import (
	"context"
	"sync"
)

type ShutdownWaiter interface {
	// Implementors of WaitUntilShutdown are expected to block until any goroutine they spawned have finished, or until the provided context has closed
	// Any persistent goroutine started with Forever is a ShutdownWaiter
	WaitUntilShutdown(shutdownContext context.Context)
}

type Supervisor interface {
	Supervise(w ShutdownWaiter)
}

type supervisedMarker interface {
	MarkSupervised()
}

// Useful for creating supervision trees; that is, nested object graphs that spawn long-running goroutines where the top level
// object needs to block until all goroutines in the systems have shut down. As such, TreeSupervisor is both a Supervisor and a ShutdownWaiter.
// When WaitUntilShutdown is called, it will in turn call WaitUntilShutdown on all of its Supervised ShutdownWaiters.
//
// Note that after calling WaitUntilShutdown it is no longer possible to call Supervise, and any subsequent call will panic.
type TreeSupervisor struct {
	supervised            []ShutdownWaiter
	waitForShutdownCalled struct {
		sync.Mutex
		called bool
	}
}

func (t *TreeSupervisor) WaitUntilShutdown(shutdownContext context.Context) {
	t.waitForShutdownCalled.Lock()
	defer t.waitForShutdownCalled.Unlock()
	t.waitForShutdownCalled.called = true
	for _, w := range t.supervised {
		w.WaitUntilShutdown(shutdownContext)
	}
}

func (t *TreeSupervisor) Supervise(w ShutdownWaiter) {
	if s, ok := w.(supervisedMarker); ok {
		s.MarkSupervised()
	}

	t.waitForShutdownCalled.Lock()
	defer t.waitForShutdownCalled.Unlock()
	if t.waitForShutdownCalled.called {
		panic("Can't call Supervise() after WaitUntilShutdown has been called")
	}
	t.supervised = append(t.supervised, w)
}
