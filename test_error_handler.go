package govnr

type report struct {
	err error
}

func (r *report) String() string {
	return r.err.Error()
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

func bufferedLogger() *collector {
	c := &collector{errors: make(chan report, 10)}
	return c
}
