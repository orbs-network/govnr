# govnr
![one does not simply start a goroutine.](https://raw.githubusercontent.com/orbs-network/govnr/master/one-does-not-simply.jpg)

[![CI](https://circleci.com/gh/orbs-network/govnr/tree/master.svg?style=svg)](https://circleci.com/gh/orbs-network/govnr/tree/master)

Use `govnr` to launch supervised goroutines. 

The package offers:
* `Once()` launches a goroutine and logs uncaught panics.
* `Forever()` launches a goroutine and in the event of a panic, log the error and re-launches, as long as the context has not been cancelled.
* `Recover()` runs a function inline, in the currently running goroutine. panics are recovered, logged and ignored.

[Docs](https://godoc.org/github.com/orbs-network/govnr) are available but could probably be better. PRs will be appreciated!

Used extensively in [Orbs Network's golang implementation](https://github.com/orbs-network/orbs-network-go) to make sure all background processes play nicely.

Example usage:
```golang

type stdoutErrorer struct {}

func (s *stdoutErrorer) Error(err error) {
	fmt.Println(err.Error())
}

errorHandler := &stdoutErrorer{}
ctx, cancel := context.WithCancel(context.Background())

data := make(chan int)
handle := govnr.Forever(ctx, "an example process", errorHandler, func() {
	for {
		select {
		case i := <-data:
			fmt.Printf("goroutine got data: %d\n", i)
		case <-ctx.Done():
			return
		}
	}
})

supervisor := &govnr.TreeSupervisor{}
supervisor.Supervise(handle)

data <- 3
data <- 2
data <- 1
cancel()

shutdownCtx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
supervisor.WaitUntilShutdown(shutdownCtx)
```
