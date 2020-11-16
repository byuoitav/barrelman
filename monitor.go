package barrelman

import "time"

// DeviceMonitor runs registered checkers on registered devices and reports on
// the outcomes of the checks according to the individual monitor's
// implementation details.
type DeviceMonitor interface {
	RegisterChecker(name string, c *Checker) error
	RegisterDevice(*Device) error

	// ForceCheck forces the DeviceMonitor to immediately run all checks for the
	// given device name.
	ForceCheck(name string) error
	Status(name string) (DeviceStatus, error)
}

type RoomMonitor interface {
}

// DeviceStatus is a representation of a device's monitoring status
type DeviceStatus struct {
	Device *Device

	// A device is typically considered "healthy" if the last run of each
	// checker passed successfully, though this can be handled differently
	// based on the DeviceMonitor
	Healthy bool

	// CheckStatus is a map of all the checkers being run by the
	// DeviceMonitor to their detailed status information
	CheckStatus map[string]DeviceCheckStatus
}

// DeviceCheckStatus represents the status of a check run on a specific device.
type DeviceCheckStatus struct {
	// RunTime is the time when the last check was run
	RunTime time.Time

	// Passed is true if the check ran successfully on the device last time
	Passed bool

	// Message is the message returned from the checker from the last run for
	// the given device
	Message string

	// Error is the error returned from the last run of the checker for this
	// device. This will be empty if passed is true in most cases
	Error string
}
