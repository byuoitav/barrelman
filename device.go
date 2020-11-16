package barrelman

// Device represents a physical device to be monitored
type Device struct {
	Name    string
	Address string

	// CheckerConfig is a map, where the key is the name of a checker, and the
	// value is arbitrary configuration values for that checker for the given device
	CheckerConfig map[string]interface{}
}

// DeviceStore is the interface to be met by the storage mechanism for device information
type DeviceStore interface {
	GetDevice(name string) (*Device, error)
	GetAllDevices() ([]*Device, error)
}
