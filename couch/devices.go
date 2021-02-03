package couch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/byuoitav/barrelman"
	"github.com/go-kivik/kivik"
)

const _centralMonitoringDB = "central-monitoring"
const _devicesDB = "shipyard-devices"

type centralDevicesDoc struct {
	Devices []centralDevice `json:"devices"`
}

type centralDevice struct {
	Name          string                 `json:"name"`
	Address       string                 `json:"address"`
	Room          string                 `json:"room"`
	CheckerConfig map[string]interface{} `json:"checkerConfig"`
}

type device struct {
	Name    string `json:"_id"`
	Address string `json:"address"`
	Room    string `json:"room"`
}

func (s *Service) GetDevice(id string) (barrelman.Device, error) {
	db := s.client.DB(context.TODO(), _devicesDB)
	dev := device{}
	err := db.Get(context.TODO(), id).ScanDoc(&dev)
	if err != nil {
		// Not found error
		if kivik.StatusCode(err) == http.StatusNotFound {
			return barrelman.Device{}, fmt.Errorf("Device not found")
		}

		return barrelman.Device{}, fmt.Errorf("couch/GetDevice get doc: %w", err)
	}

	return convertDevice(dev), nil
}

func (s *Service) GetRoomDevices(roomID string) ([]barrelman.Device, error) {
	db := s.client.DB(context.TODO(), _devicesDB)

	// Query
	q := query{
		Selector: map[string]interface{}{
			"room": roomID,
		},
		Limit: 100,
	}

	// Make the request
	rows, err := db.Find(context.TODO(), q)
	if err != nil {
		return nil, fmt.Errorf("couch/GetRoomDevices couch request: %w", err)
	}

	// Convert the devices
	devs := []barrelman.Device{}
	for rows.Next() {
		d := device{}
		err := rows.ScanDoc(&d)
		if err != nil {
			return nil, fmt.Errorf("couch/GetRoomDevices unmarshal: %w", err)
		}
		devs = append(devs, convertDevice(d))
	}

	return devs, nil
}

func (s *Service) GetCentralMonitoringDevice(name string) (*barrelman.Device, error) {
	devs, err := s.GetAllCentralMonitoringDevices()
	if err != nil {
		return nil, fmt.Errorf("get all devices: %w", err)
	}

	// Find the right device
	for _, d := range devs {
		if d.Name == name {
			return d, nil
		}
	}

	return nil, fmt.Errorf("Device not found")
}

func (s *Service) GetAllCentralMonitoringDevices() ([]*barrelman.Device, error) {
	db := s.client.DB(context.TODO(), _centralMonitoringDB)

	doc := centralDevicesDoc{}
	err := db.Get(context.TODO(), "default").ScanDoc(&doc)
	if err != nil {
		return nil, fmt.Errorf("retrieving central monitoring doc: %w", err)
	}

	// Convert types
	devs := []*barrelman.Device{}
	for _, d := range doc.Devices {
		devs = append(devs, &barrelman.Device{
			Name:    d.Name,
			Address: d.Address,
		})
	}

	return devs, nil
}

func convertDevice(d device) barrelman.Device {
	return barrelman.Device{
		Name:    d.Name,
		Address: d.Address,
		Room:    d.Room,
	}
}
