package ping

import (
	"fmt"
	"time"

	"github.com/byuoitav/barrelman"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/go-ping/ping"
)

// Checker attempts to ping the device to check for network layer health
type Checker struct {
	numPings int
	interval int
	timeout  int

	messenger *messenger.Messenger
}

// NewChecker returns a ping checker with the given options set
func NewChecker(opts ...Option) (*Checker, error) {
	c := Checker{
		numPings: 3,
		interval: 1,
		timeout:  5,
	}

	// Apply options
	for _, opt := range opts {
		opt(&c)
	}

	return &c, nil
}

// Check attempts to ping the given device. If all pings return successfully
// then the check is considered healthy. Any missed packets will result in an
// unhealthy check
func (c *Checker) Check(d *barrelman.Device, forceRecheck bool) barrelman.CheckResult {
	result := barrelman.CheckResult{
		RunTime: time.Now(),
		Passed:  true,
		Event: barrelman.Event{
			Device: d,
			Key:    "online",
			Value:  "Online",
		},
	}

	pinger, err := ping.NewPinger(d.Address)
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("Failed to creating pinger: %s", err)
		result.Event.Value = "Offline"
		return result
	}

	pinger.SetPrivileged(true)
	pinger.Count = c.numPings
	pinger.Interval = time.Duration(c.interval) * time.Second
	pinger.Timeout = time.Duration(c.timeout) * time.Second

	pinger.Run()

	stats := pinger.Statistics()

	if stats.PacketLoss == 0 {
		result.Message = fmt.Sprintf(
			"All pings returned successfully with average RTT of %fms",
			float64(stats.AvgRtt/time.Nanosecond)/1000000, // Getting ms down to several decimal places
		)
		return result
	}

	result.Error = fmt.Sprintf("Lost %d of %d pings", c.numPings-stats.PacketsRecv, c.numPings)
	result.Passed = false
	result.Event.Value = "Offline"
	return result
}
