package barrelman

// A Checker is used to "check" a specific aspect of a device for monitoring
// purposes. Examples include:
//   Ping Checker - To ensure the device is online
//   Health Checker - To check that the device is healthy
//   Serial Checker - To ensure the device hasn't been switched out
//   etc.

// Checker is the interface which will be met by the implementation
// of a monitoring "checker" to check the health of a device. A checker
// only returns a "message" as a string and an error. The message can be
// used to pass back arbitrary data from a successful check for logging
// purposes
type Checker interface {
	Check(*Device) (string, error)
}
