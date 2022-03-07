package gin

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/barrelman/localsystem"
	"github.com/gin-gonic/gin"
)

// DeviceInfo
type deviceInfo struct {
	Hostname             string `json:"hostname,omitempty"`
	ID                   string `json:"id,omitempty"`
	IP                   string `json:"ip,omitempty"`
	InternetConnectivity bool   `json:"internet-connectivity"`

	DHCPInfo struct {
		Enabled   bool `json:"enabled"`
		Togleable bool `json:"toggleable"`
	} `json:"dhcp"`
}

//getDeviceInfo responds with the list of the device info as JSON
func getDeviceInfo(c *gin.Context) {

	var info deviceInfo
	var err error

	info.Hostname, err = localsystem.Hostname()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	info.ID, err = localsystem.SystemID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	ip, err := localsystem.IPAddress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	info.IP = ip.String()
	info.InternetConnectivity = localsystem.IsConnectedToInternet()

	info.DHCPInfo.Enabled, err = localsystem.UsingDHCP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = localsystem.CanToggleDHCP()
	if err != nil {
		info.DHCPInfo.Togleable = false
	} else {
		info.DHCPInfo.Togleable = true
	}

	c.IndentedJSON(http.StatusOK, info)
}

// GetHostname returns the device's hostname
func getHostname(c *gin.Context) {
	hostname, err := localsystem.Hostname()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, hostname)
}

// GetDeviceID returns the device's id
func getDeviceID(c *gin.Context) {
	id, err := localsystem.SystemID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, id)
}

// GetIPAddress returns the device's ip address
func getIPAddress(c *gin.Context) {
	ip, err := localsystem.IPAddress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, ip.String())
}

// IsConnectedToInternet returns whether the device is connected to the internet or not
func isConnectedToInternet(c *gin.Context) {
	status := localsystem.IsConnectedToInternet()
	c.IndentedJSON(http.StatusOK, fmt.Sprintf("%v", status))
}

// GetDHCPState returns whether or not dhcp is enabled and if it can be toggled or not
func getDHCPState(c *gin.Context) {
	ret := make(map[string]interface{})

	usingDHCP, err := localsystem.UsingDHCP()
	if err != nil {
		ret["error"] = fmt.Sprintf("%v", err)
		c.JSON(http.StatusInternalServerError, ret)
		return
	}
	ret["enabled"] = usingDHCP

	if err = localsystem.CanToggleDHCP(); err != nil {
		ret["err"] = fmt.Sprintf("%v", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	ret["toggleable"] = true

	c.JSON(http.StatusOK, ret)
}
