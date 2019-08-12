/*
Package govnr provides the facilities for running supervised, persistent goroutines ("background services") that, upon unexpected panics, re-run the provided function.
This is similar to other libraries such as https://github.com/thejerf/suture, with several differences:

1. govnr accepts a context.Context and listens to its .Done() channel

2. starting a persistent process requires nothing more than calling govnr.Forever()

3. no reliance on global state
 */
package govnr
