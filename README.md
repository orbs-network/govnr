# govnr

[![CI](https://circleci.com/gh/orbs-network/govnr/tree/master.svg?style=svg)](https://circleci.com/gh/orbs-network/govnr/tree/master)

Use `govnr` to launch supervised goroutines. 

The package offers:
* `GoOnce()` launches a goroutine and ensure panics are properly logged.
* `GoForever()` launches a goroutine and in the event of a panic, log the error and relaunch until context cancellation.
* `Recover()` runs a function inline in the currently running goroutine. panics are recovered, logged and ignored.
