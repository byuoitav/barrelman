package ping

// Option is a function which modifies a given checker, allowing the
// user to have an option on how to setup the checker
type Option func(*Checker)

// WithTimeout allows the user to set the overall timeout of the ping
// check. This timeout spans all of the pings and is not specific to one
// single ping.
func WithTimeout(t int) Option {
	return func(c *Checker) {
		c.timeout = t
	}
}

// WithInterval allows the user to set the interval between pings in seconds
func WithInterval(i int) Option {
	return func(c *Checker) {
		c.interval = i
	}
}

// WithCount allows the user to set the number of pings sent to a device
// in a check
func WithCount(i int) Option {
	return func(c *Checker) {
		c.numPings = i
	}
}
