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
func (c *Checker) Check(d *barrelman.Device) (string, error) {

	pinger, err := ping.NewPinger(d.Address)
	if err != nil {
		return "", fmt.Errorf("creating pinger: %w", err)
	}

	pinger.Count = c.numPings
	pinger.Interval = time.Duration(c.interval) * time.Second
	pinger.Timeout = time.Duration(c.timeout) * time.Second

	pinger.Run()

	stats := pinger.Statistics()

	if stats.PacketLoss == 0 {
		return fmt.Sprintf(
			"All pings returned successfully with average RTT of %fms",
			float64(stats.AvgRtt/time.Millisecond),
		), nil
	}

	return "",
		fmt.Errorf("Lost %d of %d pings", c.numPings-stats.PacketsRecv, c.numPings)
}
