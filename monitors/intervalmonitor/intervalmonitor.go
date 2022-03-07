package intervalmonitor

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/byuoitav/barrelman"
)

// Monitor contains all of the data used by the IntervalMonitor
type Monitor struct {
	// Options
	jitter       int
	eventEmitter barrelman.EventEmitter

	checkers       map[string]barrelman.Checker
	checkStateChan chan deviceCheckMsg

	deviceMu sync.RWMutex
	devices  map[string]barrelman.DeviceStatus
}

type wrappedChecker struct {
	c         barrelman.Checker
	name      string
	stateChan chan deviceCheckMsg
}

type deviceCheckMsg struct {
	deviceID string
	checker  string
	result   *barrelman.CheckResult
}

// NewMonitor returns a new IntervalMonitor with the given options set
func NewMonitor(opts ...Option) (*Monitor, error) {
	m := Monitor{
		jitter:         30,
		eventEmitter:   nil,
		devices:        make(map[string]barrelman.DeviceStatus),
		checkers:       make(map[string]barrelman.Checker),
		checkStateChan: make(chan deviceCheckMsg, 100),
	}

	// Apply options
	for _, opt := range opts {
		opt(&m)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	go m.listenForChecks()

	return &m, nil
}

// listenForChecks listens for updates on device checks and writes them
// to the device state
func (m *Monitor) listenForChecks() {
	for {
		msg := <-m.checkStateChan

		// Write the new check to the device stateChan
		m.deviceMu.Lock()
		m.devices[msg.deviceID].CheckStatus[msg.checker] = *msg.result
		m.deviceMu.Unlock()

		// If there is an event emitter then send the event
		if m.eventEmitter != nil {
			go m.eventEmitter.Send(msg.result.Event)
		}
	}
}

// RegisterChecker registers the given checker under the given name to be run on
// all devices registered in this monitor on the given interval (measured in seconds)
func (m *Monitor) RegisterChecker(name string, interval int, c barrelman.Checker) error {
	// Check for existing checker
	if _, ok := m.checkers[name]; ok {
		return fmt.Errorf("Checker already registered with name %s", name)
	}

	wc := &wrappedChecker{
		c:         c,
		name:      name,
		stateChan: m.checkStateChan,
	}

	// Register checker
	m.checkers[name] = wc

	// Start checker
	go m.intervalChecker(wc, interval)

	return nil
}

// Check re-implements the barrelman.Checker interface but allows the check
// results to be sent back through the monitor's channel
func (wc *wrappedChecker) Check(d *barrelman.Device, recheck bool) barrelman.CheckResult {
	log.Printf("Running checker %s on device %s\n", wc.name, d.Name)
	result := wc.c.Check(d, recheck)

	log.Printf("Result: %+v\n", result)

	wc.stateChan <- deviceCheckMsg{
		deviceID: d.Name,
		checker:  wc.name,
		result:   &result,
	}

	return result
}

// RegisterDevice registers the given device to have all the registered checks
// run against it on an interval
func (m *Monitor) RegisterDevice(d *barrelman.Device) error {
	// Register device
	m.deviceMu.Lock()
	m.devices[d.Name] = barrelman.DeviceStatus{
		Device:      d,
		Healthy:     false,
		CheckStatus: make(map[string]barrelman.CheckResult),
	}
	m.deviceMu.Unlock()

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

// check is the internal function to run all checks on a device
func (m *Monitor) check(deviceName string) {
	// Get device
	m.deviceMu.RLock()
	d := m.devices[deviceName]
	m.deviceMu.RUnlock()

	// Run all checkers
	for _, c := range m.checkers {
		c.Check(d.Device, true)
	}
}

// Status will return the current status of the given device (by name)
func (m *Monitor) Status(name string) (barrelman.DeviceStatus, error) {
	m.deviceMu.RLock()
	status, ok := m.devices[name]
	m.deviceMu.RUnlock()

	if ok {
		return status, nil
	}

	return barrelman.DeviceStatus{}, fmt.Errorf("No device found with name %s", name)
}

// intervalChecker is the internal function used to continuously run a checker
// on all devices at the configured interval and jitter
func (m *Monitor) intervalChecker(c *wrappedChecker, interval int) {
	for {
		// Sleep for the standard interval minus a random number of seconds from
		// 0 to m.jitter. This helps to distribute large amounts of checks better
		interval := time.Duration(interval-rand.Intn(m.jitter)) * time.Second
		log.Printf("Sleeping checker %s for %s\n", c.name, interval)
		time.Sleep(interval)

		// Initiate checks on all devices
		log.Printf("Checker %s ranging devices\n", c.name)
		m.deviceMu.RLock()
		for _, d := range m.devices {
			go c.Check(d.Device, false)
		}
		m.deviceMu.RUnlock()
	}
}
