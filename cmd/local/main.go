package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/byuoitav/barrelman/avevent"
	"github.com/byuoitav/barrelman/checkers/health"
	"github.com/byuoitav/barrelman/checkers/ping"
	"github.com/byuoitav/barrelman/couch"
	"github.com/byuoitav/barrelman/monitors/intervalmonitor"
	"github.com/spf13/pflag"
)

func main() {
	var (
		dbAddr       string
		dbUser       string
		dbPass       string
		eventHubAddr string
		avAPIAddr    string
		systemID     string
	)

	pflag.StringVar(&dbAddr, "db-address", "", "The address to the couch database")
	pflag.StringVar(&dbUser, "db-username", "", "The username for the couch database")
	pflag.StringVar(&dbPass, "db-password", "", "The password for the couch database")
	pflag.StringVar(&eventHubAddr, "eventhub-address", "", "The address for the event hub")
	pflag.StringVar(&avAPIAddr, "av-api-address", "", "The address for the av api")
	pflag.StringVar(&systemID, "system-id", "", "The ID of this system")

	pflag.Parse()

	systemParts := strings.Split(systemID, "-")
	if len(systemParts) != 3 {
		log.Panicf("Invalid System ID: %s", systemID)
	}

	roomID := fmt.Sprintf("%s-%s", systemParts[0], systemParts[1])

	c, err := couch.New(dbAddr, dbUser, dbPass)
	if err != nil {
		log.Panicf("Failed to initialize couch: %s", err)
	}

	devs, err := c.GetRoomDevices(roomID)
	if err != nil {
		log.Panicf("Failed to get devices from database: %s", err)
	}

	e, err := avevent.NewLogEmitter(eventHubAddr, systemID)
	if err != nil {
		log.Panicf("Failed to start event emitter: %s", err)
	}

	m, err := intervalmonitor.NewMonitor(intervalmonitor.WithEventEmitter(e), intervalmonitor.WithJitter(5))
	if err != nil {
		log.Panicf("Failed to create interval monitor: %s", err)
	}

	pingChecker, err := ping.NewChecker()
	if err != nil {
		log.Panicf("Failed to initialize ping checker: %s", err)
	}

	healthChecker, err := health.NewChecker(avAPIAddr)
	if err != nil {
		log.Panicf("Failed to initialize health checker: %s", err)
	}

	m.RegisterChecker("ping", 120, pingChecker)
	m.RegisterChecker("health", 60, healthChecker)

	log.Printf("Beginning monitoring...")

	// Initialize monitoring
	for _, d := range devs {
		m.RegisterDevice(&d)
	}

	log.Printf("Monitoring initialized on %d devices", len(devs))

	// Hang forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
