package couch

import (
	"context"
	"fmt"

	"github.com/byuoitav/barrelman"
)

const _centralMonitoringDB = "central-monitoring"

type devicesDoc struct {
	Devices []device `json:"devices"`
}

type device struct {
	Name          string                 `json:"name"`
	Address       string                 `json:"address"`
	CheckerConfig map[string]interface{} `json:"checkerConfig"`
}

func (s *Service) GetDevice(name string) (*barrelman.Device, error) {
	devs, err := s.GetAllDevices()
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

func (s *Service) GetAllDevices() ([]*barrelman.Device, error) {
	db := s.client.DB(context.TODO(), _centralMonitoringDB)

	doc := devicesDoc{}
	err := db.Get(context.TODO(), "default").ScanDoc(&doc)
	if err != nil {
		return nil, fmt.Errorf("retrieving central monitoring doc: %w", err)
	}

	// Convert types
	devs := []*barrelman.Device{}
	for _, d := range doc.Devices {
		devs = append(devs, &barrelman.Device{
			Name:          d.Name,
			Address:       d.Address,
			CheckerConfig: d.CheckerConfig,
		})
	}

	return devs, nil
}
