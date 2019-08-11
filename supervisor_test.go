// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package govnr

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type report struct {
	err error
}

type collector struct {
	errors chan report
}

func (c *collector) Error(err error) {
	c.errors <- report{err}
}

func mockLogger() *collector {
	c := &collector{errors: make(chan report)}
	return c
}

func localFunctionThatPanics() {
	panic("foo")
}

func TestGoOnce_ReportsOnPanic(t *testing.T) {
	logger := mockLogger()

	require.NotPanicsf(t, func() {
		GoOnce(logger, localFunctionThatPanics)
	}, "GoOnce panicked unexpectedly")

	report := <-logger.errors
	require.Error(t, report.err)

}

func TestGoForever_ReportsOnPanicAndRestarts(t *testing.T) {
	numOfIterations := 10

	logger := mockLogger()
	ctx, cancel := context.WithCancel(context.Background())

	count := 0

	require.NotPanicsf(t, func() {
		GoForever(ctx, logger, func() {
			if count > numOfIterations {
				cancel()
			} else {
				count++
			}
			panic("foo")
		})
	}, "GoForever panicked unexpectedly")

	for i := 0; i < numOfIterations; i++ {
		select {
		case report := <-logger.errors:
			require.Error(t, report.err)
		case <-time.After(1 * time.Second):
			require.Fail(t, "long living goroutine didn't restart")
		}
	}
}

func TestGoForever_TerminatesWhenContextIsClosed(t *testing.T) {
	logger := mockLogger()
	ctx, cancel := context.WithCancel(context.Background())

	bgStarted := make(chan struct{})
	bgEnded := make(chan struct{})
	shutdown := GoForever(ctx, logger, func() {
		bgStarted <- struct{}{}
		select {
		case <-ctx.Done():
			bgEnded <- struct{}{}
			return
		}
	})

	<-bgStarted
	cancel()

	select {
	case <-bgEnded:
		// ok, invocation of cancel() caused goroutine to stop, we can now check if it restarts
	case <-time.After(1 * time.Second):
		require.Fail(t, "long living goroutine didn't stop")
	}

	select {
	case <-shutdown:
		// system has shutdown, all ok
	case <-time.After(1 * time.Second):
		t.Fatalf("long living goroutine did not return")
	}

}
