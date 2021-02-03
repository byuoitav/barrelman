package barrelman

// Device represents a physical device to be monitored
type Device struct {
	Name    string
	Address string
	Room    string
}

// DeviceStore is the interface to be met by the storage mechanism for device information
type DeviceStore interface {
	GetCentralMonitoringDevice(name string) (*Device, error)
	GetAllCentralMonitoringDevices() ([]*Device, error)
	GetRoomDevices(roomID string) ([]*Device, error)
	GetDevice(name string) (*Device, error)
}
