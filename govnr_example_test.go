package govnr

import (
	"context"
	"fmt"
	"time"
)

type stdoutErrorer struct {}

func (s *stdoutErrorer) Error(err error) {
	fmt.Println(err.Error())
}

func Example() {
	errorHandler := &stdoutErrorer{}
	ctx, cancel := context.WithCancel(context.Background())

	data := make(chan int)
	handle := Forever(ctx, "an example process", errorHandler, func() {
		for {
			select {
			case i := <-data:
				fmt.Printf("goroutine got data: %d\n", i)
			case <-ctx.Done():
				return
			}
		}
	})

	supervisor := &TreeSupervisor{}
	supervisor.Supervise(handle)

	data <- 3
	data <- 2
	data <- 1
	cancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	supervisor.WaitUntilShutdown(shutdownCtx)

	// Output:
	// goroutine got data: 3
	// goroutine got data: 2
	// goroutine got data: 1
}

