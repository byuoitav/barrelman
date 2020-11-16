package main

import (
	"log"
	"sync"

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
	)

	pflag.StringVar(&dbAddr, "db-address", "", "The address to the couch database")
	pflag.StringVar(&dbUser, "db-username", "", "The username for the couch database")
	pflag.StringVar(&dbPass, "db-password", "", "The password for the couch database")
	pflag.StringVar(&eventHubAddr, "eventhub-address", "", "The address for the event hub")

	pflag.Parse()

	c, err := couch.New(dbAddr, dbUser, dbPass)
	if err != nil {
		log.Panicf("Failed to initialize couch: %s", err)
	}

	devs, err := c.GetAllDevices()
	if err != nil {
		log.Panicf("Failed to get receivers from database: %s", err)
	}

	m, err := intervalmonitor.NewMonitor()
	if err != nil {
		log.Panicf("Failed to create interval monitor: %s", err)
	}

	pingChecker, err := ping.NewChecker()
	if err != nil {
		log.Panicf("Failed to initialize ping checker: %s", err)
	}

	m.RegisterChecker("ping", pingChecker)

	log.Printf("Beginning monitoring...")

	// Initialize monitoring
	for _, d := range devs {
		m.RegisterDevice(d)
	}

	log.Printf("Monitoring initialized on %d devices", len(devs))

	// Hang forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
