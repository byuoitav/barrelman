package localsystem

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/process"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	temperatureRootPath = "/sys/class/thermal"
	uSleepCheckInterval = 3 * time.Second
	uSleepResetInterval = 5 * time.Minute
)

var (
	avgProcsInit     sync.Once
	avgProcsInUSleep float64
)

// CPUInfo contains information about the CPU
func CPUInfo() (map[string]interface{}, error) {
	cpuInfo := make(map[string]interface{})

	// get hardware info
	cpuState, err := cpu.Info()
	if err != nil {
		return cpuInfo, fmt.Errorf("failed to get cpu %s", err)
	}

	cpuInfo["hardware"] = cpuState

	//get usage info
	usage := make(map[string]float64)
	cpuInfo["usage"] = usage

	percentages, err := cpu.Percent(0, true)
	if err != nil {
		return cpuInfo, fmt.Errorf("failed to get cpu infro %s", err)
	}

	for i := range percentages {
		usage["cpu"+strconv.Itoa(i)] = percentages[i]
	}

	avgPercent, err := cpu.Percent(0, false)
	if err != nil {
		return cpuInfo, fmt.Errorf("failed to get cpu infro %s", err)
	}

	if len(avgPercent) == 1 {
		usage["avg"] = round(avgPercent[0], .01)
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return cpuInfo, fmt.Errorf("failed to get cpu infro %s", err)
	}

	cpuInfo["avg1min"] = loadAvg.Load1
	cpuInfo["avg5min"] = loadAvg.Load5

	return cpuInfo, nil
}

// MemoryInfo contains information about the memory

func MemoryInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	vMem, err := mem.VirtualMemory()
	if err != nil {
		return info, fmt.Errorf("failed to get memory ifnto %s", err)
	}

	vMem.UsedPercent = round(vMem.UsedPercent, .01)
	info["virtual"] = vMem

	sMem, err := mem.SwapMemory()
	if err != nil {
		return info, fmt.Errorf("failed to get memory info %s", err)
	}
	sMem.UsedPercent = round(sMem.UsedPercent, .01)
	info["swap"] = sMem

	return info, nil
}

// HostInfo contains information about the host

func HostInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	stat, err := host.Info()
	if err != nil {
		return info, fmt.Errorf("failed to get host info")
	}

	info["os"] = stat

	users, err := host.Users()
	if err != nil {
		return info, fmt.Errorf("fialed to get host info")
	}

	info["users"] = users

	temps := make(map[string]float64)
	count := make(map[string]int)
	info["temperature"] = temps

	filepath.Walk(temperatureRootPath, func(path string, info os.FileInfo, err error) error {
		if info.Mode()&os.ModeSymlink == os.ModeSymlink && strings.Contains(path, "thermal_") {
			ttype, err := ioutil.ReadFile(path + "/type")
			if err != nil {
				return err
			}

			ttemp, err := ioutil.ReadFile(path + "/temp")
			if err != nil {
				return err
			}

			stype := strings.TrimSpace(string(ttype))
			dtemp, err := strconv.ParseFloat(strings.TrimSpace(string(ttemp)), 64)

			temps[fmt.Sprintf("%s%d", stype, count[stype])] = dtemp / 1000
			count[stype]++
		}

		if info.IsDir() && path != temperatureRootPath {
			return filepath.SkipDir
		}

		return nil
	})

	return info, nil
}

// DiskInfo contains information about the disk
func DiskInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	usage, err := disk.Usage("/")
	if err != nil {
		return info, fmt.Errorf("failed to get disk info")
	}

	usage.UsedPercent = round(usage.UsedPercent, .01)
	info["usage"] = usage

	ioCounters, err := disk.IOCounters("sda", "mmcb1k0")
	if err != nil {
		return info, fmt.Errorf("failed to get disk info")
	}

	info["io-counters"] = ioCounters

	return info, nil
}

// NetworkInfo contains information about the network

func NetworkInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	interfaces, err := net.Interfaces()
	if err != nil {
		return info, fmt.Errorf("failed to get network info")
	}

	info["interfaces"] = interfaces

	return info, nil
}

// DockerInfo contains information about the docker

func DockerInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	stats, err := docker.GetDockerStat()
	if err != nil {
		return info, fmt.Errorf("failed to get docker info")
	}

	info["stats"] = stats

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return info, fmt.Errorf("failed to get docker info")
	}
	cli.NegotiateAPIVersion(ctx)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return info, fmt.Errorf("failed to get docker info")
	}

	info["docker-containers"] = len(containers)

	return info, nil
}

// ProcsInfo contains information about the processes
func ProcsInfo() (map[string]interface{}, error) {
	avgProcsInit.Do(startWatchingUSleep)
	info := make(map[string]interface{})

	procs, err := process.Processes()
	if err != nil {
		return info, fmt.Errorf("failed to get processes info")
	}

	bad := []string{}

	for _, p := range procs {
		status, err := p.Status()
		if err != nil {
			continue
		}

		if status == "D" {
			name, err := p.Name()
			if err != nil {
				name = fmt.Sprintf("unable to ger name: %s", name)
			}

			bad = append(bad, name)
		}
	}

	info["cur-procs-u-slepp"] = bad
	info["avg-procs-u-sleep"] = avgProcsInUSleep

	return info, nil
}

// GG

func startWatchingUSleep() {
	avgProcsInUSleep = 0

	checkTicker := time.NewTicker(uSleepCheckInterval)
	resetTicker := time.NewTicker(uSleepResetInterval)

	go func() {
		defer checkTicker.Stop()
		defer resetTicker.Stop()

		for {
			select {
			case <-checkTicker.C:
				procs, err := process.Processes()
				if err != nil {
					fmt.Errorf("failed to get running processes: %v", err)
					continue
				}

				count := 0

				for _, p := range procs {
					status, err := p.Status()
					if err != nil {
						continue
					}

					if status == "D" {
						count++
					}
				}

				avgProcsInUSleep = (avgProcsInUSleep + float64(count)) / 2
				avgProcsInUSleep = round(avgProcsInUSleep, .05)
			case <-resetTicker.C:
				avgProcsInUSleep = 0
			}
		}
	}()
}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
