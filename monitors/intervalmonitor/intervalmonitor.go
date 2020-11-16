package intervalmonitor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/byuoitav/barrelman"
)

// Monitor contains all of the data used by the IntervalMonitor
type Monitor struct {
	// Options
	interval int
	jitter   int

	devices  map[string]barrelman.DeviceStatus
	checkers map[string]barrelman.Checker
}

// NewMonitor returns a new IntervalMonitor with the given options set
func NewMonitor(opts ...Option) (*Monitor, error) {
	m := Monitor{
		interval: 600,
		jitter:   30,
		devices:  make(map[string]barrelman.DeviceStatus),
		checkers: make(map[string]barrelman.Checker),
	}

	// Apply options
	for _, opt := range opts {
		opt(&m)
	}

	return &m, nil
}

// RegisterChecker registers the given checker under the given name to be run on
// all devices registered in this monitor
func (m *Monitor) RegisterChecker(name string, c barrelman.Checker) error {
	// Check for existing checker
	if _, ok := m.checkers[name]; ok {
		return fmt.Errorf("Checker already registered with name %s", name)
	}

	// Register checker
	m.checkers[name] = c
	return nil
}

// RegisterDevice registers the given device to have all the registered checks
// run against it on an interval
func (m *Monitor) RegisterDevice(d *barrelman.Device) error {
	// Check for already existing device
	if _, ok := m.devices[d.Name]; ok {
		return fmt.Errorf("Device already registered with name %s", d.Name)
	}

	// Register device
	m.devices[d.Name] = barrelman.DeviceStatus{
		Device:      d,
		Healthy:     false,
		CheckStatus: make(map[string]barrelman.DeviceCheckStatus),
	}

	// Start periodic checks
	go m.intervalCheck(d.Name)

	return nil
}

// ForceCheck forces the monitor to immediately run all registered checkers against
// the previously registered device by its name
func (m *Monitor) ForceCheck(name string) error {
	if _, ok := m.devices[name]; ok {
		m.check(name)
		return nil
	}

	return fmt.Errorf("No device found with name %s", name)
}

// Status will return the current status of the given device (by name)
func (m *Monitor) Status(name string) (barrelman.DeviceStatus, error) {
	if status, ok := m.devices[name]; ok {
		return status, nil
	}

	return barrelman.DeviceStatus{}, fmt.Errorf("No device found with name %s", name)
}

// intervalCheck is the internal function used to continuously check on a device
// at the configured interval and jitter
func (m *Monitor) intervalCheck(deviceName string) {
	for {
		m.check(deviceName)

		// Sleep for the standard interval minus a random number of seconds from
		// 0 to m.jitter. This helps to distribute large amounts of checks better
		interval := time.Duration(m.interval-rand.Intn(m.jitter)) * time.Second
		time.Sleep(interval)
	}
}

// check is the internal function to run all checks on a device
func (m *Monitor) check(deviceName string) {
	for checkName, c := range m.checkers {
		status := barrelman.DeviceCheckStatus{
			RunTime: time.Now(),
		}

		msg, err := c.Check(m.devices[deviceName].Device)
		if err != nil {
			status.Error = err.Error()
			status.Passed = false
		} else {
			status.Message = msg
			status.Passed = true
		}

		m.devices[deviceName].CheckStatus[checkName] = status
	}
}
