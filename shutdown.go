package govnr

import (
	"context"
	"sync"
)

type ShutdownWaiter interface {
	WaitUntilShutdown(shutdownContext context.Context)
}

type Supervisor interface {
	MarkSupervised()
}

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
	if s, ok := w.(Supervisor); ok {
		s.MarkSupervised()
	}

	t.waitForShutdownCalled.Lock()
	defer t.waitForShutdownCalled.Unlock()
	if t.waitForShutdownCalled.called {
		panic("Can't call Supervise() after WaitUntilShutdown has been called")
	}
	t.supervised = append(t.supervised, w)
}
