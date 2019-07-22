# govnr

[![CI](https://circleci.com/gh/orbs-network/govnr/tree/master.svg?style=svg)](https://circleci.com/gh/orbs-network/govnr/tree/master)

Use `govnr` to launch supervised go routines. 

The package offers:
* `GoOnce()   ` - launch a goroutine and ensure panics are properly logged, a
* `GoForever()` - launch a goroutine and in the event of a panic, log the error and relaunch until context cancellation.
* `Recover()  ` - a helper method to run a function inline in the currently running goroutine. panics are recovered, logged and ignored.
