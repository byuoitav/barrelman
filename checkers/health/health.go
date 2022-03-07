package health

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/byuoitav/barrelman"
	"golang.org/x/sync/singleflight"
)

// roomHealth and deviceHealth encapsulate the response from the av api
type roomHealth struct {
	Devices map[string]deviceHealth `json:"devices"`
	Expires time.Time
}

type deviceHealth struct {
	Healthy *bool   `json:"healthy"`
	Error   *string `json:"error"`
}

// Option is an option for the checker
type Option func(*Checker)

// Checker is a health checker for devices
type Checker struct {
	apiAddress   string
	cacheTimeout int
	sfGroup      singleflight.Group

	cacheMu sync.RWMutex
	cache   map[string]roomHealth
}

// NewChecker returns a new Health Checker which will hit the health endpoint
// on the given API Address given the options passed in
func NewChecker(apiAddress string, opts ...Option) (*Checker, error) {
	return &Checker{
		apiAddress:   apiAddress,
		cacheTimeout: 45, // 45 seconds by default
		cache:        make(map[string]roomHealth),
	}, nil
}

// Check will hit the AV Control API health endpoint for the given device.
// If the response is a 200 then the device is considered healthy. All other
// status codes will return an unhealthy check
func (c *Checker) Check(d *barrelman.Device, recheck bool) barrelman.CheckResult {
	result := barrelman.CheckResult{
		RunTime: time.Now(),
		Passed:  true,
		Event: barrelman.Event{
			Device: d,
			Key:    "responsive",
			Value:  "Ok",
		},
	}

	// Get cache entry for room
	c.cacheMu.RLock()
	rHealth, ok := c.cache[d.Room]
	c.cacheMu.RUnlock()

	// If we are forcing a recheck, the cache doesn't have an entry,
	// or that entry is expired, then refresh the cache
	if recheck || !ok || (ok && rHealth.Expires.Before(time.Now())) {
		res, err, _ := c.sfGroup.Do(d.Room, func() (interface{}, error) {
			return c.refreshRoomHealth(d.Room)
		})
		if err != nil {
			result.Error = err.Error()
			result.Passed = false
			result.Event.Value = "No Response"
			return result
		}
		rHealth = res.(roomHealth)
	}

	// If the device exists in the roomHealth
	if devHealth, ok := rHealth.Devices[d.Name]; ok {
		// If the device had a health check
		if devHealth.Healthy != nil {
			// If the device is not healthy
			if !*devHealth.Healthy {
				result.Passed = false
				result.Event.Value = "No Response"
				// Try to get an error
				if devHealth.Error != nil {
					result.Error = *devHealth.Error
				} else {
					result.Error = "Unhealthy"
				}
				return result
			}
		} else { // Device didn't have a health check
			result.Message = "No health check implemented"
			return result
		}
	} else { // Device wasn't found in health response
		result.Message = "Device not found in room"
		return result
	}

	// Device is healthy
	result.Message = "Healthy"
	return result

}

func (c *Checker) refreshRoomHealth(room string) (roomHealth, error) {
	url := fmt.Sprintf("%s/api/v1/room/%s/health", c.apiAddress, room)

	// Make health request
	res, err := http.Get(url)
	if err != nil {
		return roomHealth{}, fmt.Errorf("failed to get room health: %w", err)
	}
	defer res.Body.Close()

	// Read body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return roomHealth{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Error if we didn't get a 2XX response code
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return roomHealth{}, fmt.Errorf("got non 200 status back from api: %s", body)
	}

	// Parse response
	h := roomHealth{}
	err = json.Unmarshal(body, &h)
	if err != nil {
		return roomHealth{}, fmt.Errorf("failed to parse response body: %w", err)

	}

	// Write to cache
	h.Expires = time.Now().Add(time.Duration(c.cacheTimeout) * time.Second)
	c.cacheMu.Lock()
	c.cache[room] = h
	c.cacheMu.Unlock()

	return h, nil
}
