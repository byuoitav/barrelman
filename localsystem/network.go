package localsystem

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	dhcpFile = "/etc/dhcpcd.conf"
)

// Hostname returns the hostname of the device
func Hostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname %s", err)
	}

	return hostname, nil
}

// MustHostname returns the hostname of the device, and panics if it fails
func MustHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("failed to get hostname %s", err)
	}

	return hostname
}

// IPAddress gets the public ip address of the device
func IPAddress() (net.IP, error) {
	var ip net.IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get ip address of device %s", err)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && strings.Contains(address.String(), "/24") {
			ip, _, err = net.ParseCIDR(address.String())
			if err != nil {
				return nil, fmt.Errorf("failed to get ip address of device")
			}
		}
	}

	if ip == nil {
		return nil, fmt.Errorf("failed to get ip address of device")
	}

	fmt.Printf("My IP address is %v", ip.String())
	return ip, nil
}

// IsConnectedToInternet returns wether the device can reach google's servers.
func IsConnectedToInternet() bool {
	_, err := net.Dial("tcp", "google.com:80")

	return err == nil
}

// UsingDHCP returns wether the device is using DHCP
func UsingDHCP() (bool, error) {
	// read dhcp.cong file
	contents, err := ioutil.ReadFile(dhcpFile)
	if err != nil {
		return false, fmt.Errorf("uanble to read %s", dhcpFile)
	}

	reg := regexp.MustCompile(`(?m)^static ip_address`)
	matches := reg.Match(contents)

	return !matches, nil
}

// ToggleDHCP turns dhcp on/off
func ToggleDHCP() error {
	if err := CanToggleDHCP(); err != nil {
		return err
	}
	tmpFile := fmt.Sprintf("%s.tmp", dhcpFile)
	otherFile := fmt.Sprintf("%s.other", dhcpFile)

	// swap the files
	err := os.Rename(dhcpFile, tmpFile)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	err = os.Rename(otherFile, dhcpFile)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	err = os.Rename(tmpFile, otherFile)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	_, err = exec.Command("sh", "-c", "sudo systemctl restart dhcpcd").Output()
	if err != nil {
		return fmt.Errorf("unable to restart dhcpcd service")
	}

	return nil
}

// CanToggleDHCP returns nil if you can toggle DHCP, or an errosr if yo can't
func CanToggleDHCP() error {
	otherFile := fmt.Sprintf("%s.other", dhcpFile)

	if _, err := os.Stat(dhcpFile); os.IsNotExist(err) {
		return fmt.Errorf("can't toggle dhcp because there is no %s file", dhcpFile)
	}
	if _, err := os.Stat(otherFile); os.IsNotExist(err) {
		return fmt.Errorf("can't toggle dhcp because there is no %s.other file", dhcpFile)
	}

	return nil
}
