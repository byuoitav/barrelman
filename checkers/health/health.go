package health

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/barrelman"
)

// Option is an option for the checker
type Option func(*Checker)

// Checker is a health checker for devices
type Checker struct {
	apiAddress string
}

// NewChecker returns a new Health Checker which will hit the health endpoint
// on the given API Address given the options passed in
func NewChecker(apiAddress string, opts ...Option) (*Checker, error) {
	return &Checker{
		apiAddress: apiAddress,
	}, nil
}

// Check will hit the AV Control API health endpoint for the given device.
// If the response is a 200 then the device is considered healthy. All other
// status codes will return an unhealthy check
func (c *Checker) Check(d *barrelman.Device) (string, error) {

	res, err := http.Get(fmt.Sprintf("%s/%s/health", c.apiAddress, d.Address))
	if err != nil {
		return "",
			fmt.Errorf("call health endpoint: %w", err)
	}
	defer res.Body.Close()

	// Ignore error, if we get 200 then we still consider things to be healthy
	body, _ := ioutil.ReadAll(res.Body)

	// Check for a 200
	if res.StatusCode == http.StatusOK {
		return fmt.Sprintf("Response: %s", string(body)), nil
	}

	return "", fmt.Errorf("Status Code: %d Body: %s", res.StatusCode, string(body))
}
