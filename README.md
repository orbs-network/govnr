# govnr
![alt text](https://raw.githubusercontent.com/orbs-network/govnr/master/one-does-not-simply.jpg)


[![CI](https://circleci.com/gh/orbs-network/govnr/tree/master.svg?style=svg)](https://circleci.com/gh/orbs-network/govnr/tree/master)

Use `govnr` to launch supervised goroutines. 

The package offers:
* `Once()` launches a goroutine and logs uncaught panics.
* `Forever()` launches a goroutine and in the event of a panic, log the error and re-launch if the context has not been cancelled.
* `Recover()` runs a function inline, in the currently running goroutine. panics are recovered, logged and ignored.
