package barrelman

import "time"

// A Checker is used to "check" a specific aspect of a device for monitoring
// purposes. Examples include:
//   Ping Checker - To ensure the device is online
//   Health Checker - To check that the device is healthy
//   Serial Checker - To ensure the device hasn't been switched out
//   etc.

// Checker is the interface which will be met by the implementation
// of a monitoring "checker" to check the health of a device. The checker
// returns a CheckResult which contains status information about the check
// that was run.
// The forceRecheck parameter can be set to true to force the checker
// to recheck the device rather than returning a value from cache
type Checker interface {
	Check(device *Device, forceRecheck bool) CheckResult
}

// DeviceCheckStatus represents the status of a check run on a specific device.
type CheckResult struct {
	// RunTime is the time when the last check was run
	RunTime time.Time

	// Passed is true if the check ran successfully and should be considered "passing"
	Passed bool

	// Message is an arbitrary message returned from the checker
	Message string

	// Error is an error message from the checker, this is typically null if the
	// check is considered "passed"
	Error string

	// Event is the event that should be emitted (if an emitter is used)
	// for the check
	Event Event
}
