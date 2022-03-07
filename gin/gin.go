package gin

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Service struct {
	WebRoot string
}

func (s *Service) Serve(addr string) error {
	router := gin.Default()

	// API Endpoints
	apiGroup := router.Group("/api/v1")
	apiGroup.GET("/device", getDeviceInfo)
	apiGroup.GET("/device/hostname", getHostname)
	apiGroup.GET("/device/id", getDeviceID)
	apiGroup.GET("/device/ip", getIPAddress)
	apiGroup.GET("/device/network", isConnectedToInternet)
	apiGroup.GET("/device/dhcp", getDHCPState)

	router.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)

		if file == "" || filepath.Ext(file) == "" {
			c.File(fmt.Sprintf("%s/index.html", s.WebRoot))
		} else {
			c.File(fmt.Sprintf("%s/", s.WebRoot) + path.Join(dir, file))
		}
	})

	return router.Run(addr)
}
