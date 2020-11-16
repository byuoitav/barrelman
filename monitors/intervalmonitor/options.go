package intervalmonitor

// Option is a function which modifies a Monitor. This allows the user to set
// options that have been exposed
type Option func(*Monitor)

// WithInterval allows the user to set the interval of the checks in seconds.
// The default is 600 seconds.
func WithInterval(i int) Option {
	return func(m *Monitor) {
		m.interval = i
	}
}

// WithJitter allows the user to set the amount of jitter that is applied to the
// timing of the checks in seconds. The default is 30 seconds
func WithJitter(j int) Option {
	return func(m *Monitor) {
		m.jitter = j
	}
}
