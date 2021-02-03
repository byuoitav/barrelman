package avevent

import (
	"fmt"
	"log"
	"time"

	"github.com/byuoitav/barrelman"
	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/v2/events"
)

type LogEventEmitter struct {
	m        *messenger.Messenger
	systemID string
}

func NewLogEmitter(hubAddress, systemID string) (*LogEventEmitter, error) {
	m, err := messenger.BuildMessenger(hubAddress, base.Messenger, 1000)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to build messenger: %s", err)
	}

	return &LogEventEmitter{
		m:        m,
		systemID: systemID,
	}, nil
}

func (e *LogEventEmitter) Send(event barrelman.Event) {
	// Log first
	log.Printf("Event: Key: %s | Value: %s | Device: %s", event.Key, event.Value, event.Device.Name)

	// Emit event to av central hub
	devInfo := events.GenerateBasicDeviceInfo(event.Device.Name)
	newEvent := events.Event{
		GeneratingSystem: e.systemID,
		Timestamp:        time.Now(),
		TargetDevice:     devInfo,
		AffectedRoom:     devInfo.BasicRoomInfo,
		Key:              event.Key,
		Value:            event.Value,
	}

	e.m.SendEvent(newEvent)
}
