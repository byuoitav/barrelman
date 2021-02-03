package intervalmonitor

import "github.com/byuoitav/barrelman"

// Option is a function which modifies a Monitor. This allows the user to set
// options that have been exposed
type Option func(*Monitor)

// WithJitter allows the user to set the amount of jitter that is applied to the
// timing of the checks in seconds. The default is 30 seconds
func WithJitter(j int) Option {
	return func(m *Monitor) {
		m.jitter = j
	}
}

// WithEventEmitter allows the user to set an EventEmitter to be used by the
// monitor as it runs checks
func WithEventEmitter(e barrelman.EventEmitter) Option {
	return func(m *Monitor) {
		m.eventEmitter = e
	}
}
