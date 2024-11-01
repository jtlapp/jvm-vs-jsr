package util

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type ResourceStatus struct {
	TotalAvailablePorts   uint
	TotalFileDescriptors  uint
	FDsInUseCount         uint
	TimeWaitPortsCount    uint
	EstablishedPortsCount uint
}

var portRangeSize = getPortRangeSize()
var totalFileDescriptors = getTotalFileDescriptors()

func NewResourceStatus() ResourceStatus {
	timeWaitPortsCount, establishedPortsCount := getPortsInUseCounts()
	return ResourceStatus{
		TotalAvailablePorts:   portRangeSize,
		TotalFileDescriptors:  totalFileDescriptors,
		FDsInUseCount:         getFDsInUseCount(),
		TimeWaitPortsCount:    timeWaitPortsCount,
		EstablishedPortsCount: establishedPortsCount,
	}
}

func (rs *ResourceStatus) GetPercentages() (float64, float64, float64) {
	return float64(rs.EstablishedPortsCount) * 100 / float64(rs.TotalAvailablePorts),
		float64(rs.TimeWaitPortsCount) * 100 / float64(rs.TotalAvailablePorts),
		float64(rs.FDsInUseCount) * 100 / float64(rs.TotalFileDescriptors)
}

func WaitForPortsToClear() {
	timeWaitPortsCount, establishedPortsCount := getPortsInUseCounts()
	for timeWaitPortsCount+establishedPortsCount > 0 {
		time.Sleep(time.Second)
		timeWaitPortsCount, establishedPortsCount = getPortsInUseCounts()
	}
}

func getFDsInUseCount() uint {
	inUseFDs, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		panic(err)
	}
	return uint(len(inUseFDs))
}

func getPortsInUseCounts() (timeWaitCount, establishedCount uint) {
	data, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		panic(err)
	}

	// Skip header line
	lines := strings.Split(string(data), "\n")[1:]

	// Column 4 contains the connection state in hex
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			state := fields[3]
			switch state {
			case "06":
				timeWaitCount++
			case "01":
				establishedCount++
			}
		}
	}
	return timeWaitCount, establishedCount
}

func getPortRangeSize() uint {
	data, err := os.ReadFile("/proc/sys/net/ipv4/ip_local_port_range")
	if err != nil {
		panic(err)
	}

	var lowPort, highPort int
	_, err = fmt.Sscanf(string(data), "%d %d", &lowPort, &highPort)
	if err != nil {
		panic(err)
	}
	return uint(highPort - lowPort)
}

func getTotalFileDescriptors() uint {
	var rlimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	return uint(rlimit.Cur)
}
