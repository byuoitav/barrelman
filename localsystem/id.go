package localsystem

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	EnvironmentFile = "/avcapi/environment"
)

var (
	systemID    = os.Getenv("SYSTEM_ID")
	installerID = os.Getenv("INSTALLER_ID")
)

func SystemID() (string, error) {
	if len(systemID) == 0 {
		return "", fmt.Errorf("SYSTEM_ID not set")
	}

	if IsDeviceIDValid(systemID) {
		return "", fmt.Errorf("SYSTEM_ID is set as %s, wich is an invalid hostname", systemID)
	}

	return systemID, nil
}

// RoomID returns the room ID of the pi based on the hostname
func RoomID() (string, error) {
	id, err := SystemID()
	if err != nil {
		return "", fmt.Errorf("failed to get RoomID %v", err)
	}

	split := strings.Split(id, "-")
	return split[0] + "-" + split[1], nil
}

func IsDeviceIDValid(id string) bool {
	reg := regexp.MustCompile(`([A-z,0-9]{2,}-[A-z,0-9]+)-[A-z]+[0-9]+`)
	vals := reg.FindStringSubmatch(id)
	return len(vals) != 0
}
