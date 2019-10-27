package govnr

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTreeSupervisor_SuperviseAfterWaitForShutdown_Panics(t *testing.T) {
	s := TreeSupervisor{}
	s.WaitUntilShutdown(context.Background())
	require.Panics(t, func() {
		s.Supervise(nil)
	})
}

func TestTreeSupervisor_Supervise(t *testing.T) {
	logger := bufferedLogger()
	ctx, cancel := context.WithCancel(context.Background())
	s := TreeSupervisor{}
	s.Supervise(Forever(ctx, "foo", logger, func() {
		<-ctx.Done()
	}))
	cancel()
	s.WaitUntilShutdown(context.Background())
	require.Empty(t, logger.errors)
}


