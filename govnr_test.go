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

func localFunctionThatPanics() {
	panic("foo")
}

func TestOnce_ReportsOnPanic(t *testing.T) {
	logger := mockLogger()

	require.NotPanicsf(t, func() {
		Once(logger, localFunctionThatPanics)
	}, "GoOnce panicked unexpectedly")

	report := <-logger.errors
	require.Error(t, report.err)

}

func TestForever_ReportsOnPanicAndRestarts(t *testing.T) {
	numOfIterations := 10

	logger := mockLogger()
	ctx, cancel := context.WithCancel(context.Background())

	count := 0

	require.NotPanicsf(t, func() {
		handle := Forever(ctx, "some service", logger, func() {
			if count > numOfIterations {
				cancel()
			} else {
				count++
			}
			panic("foo")
		})
		handle.MarkSupervised()
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

func TestForever_TerminatesWhenContextIsClosed(t *testing.T) {
	logger := bufferedLogger()
	ctx, cancel := context.WithCancel(context.Background())

	bgStarted := make(chan struct{})
	bgEnded := make(chan struct{})
	handle := Forever(ctx, "another service", logger, func() {
		bgStarted <- struct{}{}
		select {
		case <-ctx.Done():
			bgEnded <- struct{}{}
			return
		}
	})
	handle.MarkSupervised()

	<-bgStarted
	cancel()

	select {
	case <-bgEnded:
		// ok, invocation of cancel() caused goroutine to stop, we can now check if it restarts
	case <-time.After(1 * time.Second):
		require.Fail(t, "long living goroutine didn't stop")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	handle.WaitUntilShutdown(shutdownCtx)
	require.Empty(t, logger.errors, "error was reported on shutdown")
}

func TestForeverHandle_ErrorsWhenTerminatedWithoutSupervision(t *testing.T) {
	logger := bufferedLogger()

	h := &ForeverHandle{closed: make(chan struct{}), supervisedChan: make(chan struct{}), errorHandler: logger, name: "foo"}
	h.terminated()
	select {
	case report := <-logger.errors:
		require.EqualError(t, report.err, "Forever governed goroutine foo terminated without being supervised")
	default:
		t.Errorf("handle didn't error on termination")
	}
}

func TestForeverHandle_DoesNotErrorWhenTerminatedAfterSupervision(t *testing.T) {
	logger := bufferedLogger()

	h := &ForeverHandle{closed: make(chan struct{}), supervisedChan: make(chan struct{}), errorHandler: logger, name: "foo"}
	h.MarkSupervised()
	h.terminated()
	require.Empty(t, logger.errors, "error was reported on shutdown")
}

func TestForeverHandle_DoesNotRaceWhenContextClosedBeforeSupervision(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	logger := bufferedLogger()

	handle := Forever(ctx, "another service", logger, func() {
		t.Errorf("job should not be called")
	})
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	time.Sleep(50 * time.Millisecond)
	handle.MarkSupervised()
	handle.WaitUntilShutdown(shutdownCtx)

	require.Empty(t, logger.errors, "error was reported on shutdown")
}
