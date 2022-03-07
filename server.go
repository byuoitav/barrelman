package barrelman

import (
	"github.com/byuoitav/barrelman/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//device info endpoints
	router.GET("/device", handlers.GetDeviceInfo)
	router.GET("/device/hostname", handlers.GetHostname)
	router.GET("/device/id", handlers.GetDeviceID)
	router.GET("/device/ip", handlers.GetIPAddress)
	router.GET("/device/network", handlers.IsConnectedToInternet)
	router.GET("/device/dhcp", handlers.GetDHCPState)

	router.Run("localhost:8080")
}
