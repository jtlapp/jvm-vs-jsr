package util

import (
	"fmt"
	"runtime"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/platform"
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

func PortsAreReady(maxReservedPorts uint) (bool, error) {
	timeWaitPortsCount, establishedPortsCount := getPortsInUseCounts()
	if establishedPortsCount > maxReservedPorts {
		return false, fmt.Errorf(
			"expected at most %d active ports but found %d (boost %s if this is okay)",
			maxReservedPorts, establishedPortsCount, config.MaxReservedPortsEnvVar)
	}
	return timeWaitPortsCount == 0, nil
}

func WaitForPortsToTimeout() {
	timeWaitPortsCount, _ := getPortsInUseCounts()
	for timeWaitPortsCount > 0 {
		time.Sleep(time.Second)
		timeWaitPortsCount, _ = getPortsInUseCounts()
	}
}

func getFDsInUseCount() uint {
	switch runtime.GOOS {
	case "darwin":
		return platform.GetFDsInUseCountOnMac()
	case "windows":
		return platform.GetFDsInUseCountOnWindows()
	default:
		return platform.GetFDsInUseCountOnLinux()
	}
}

func getPortsInUseCounts() (timeWaitCount, establishedCount uint) {
	switch runtime.GOOS {
	case "darwin":
		return platform.GetPortsInUseCountsOnMac()
	case "windows":
		return platform.GetPortsInUseCountsOnWindows()
	default:
		return platform.GetPortsInUseCountsOnLinux()
	}
}

func getPortRangeSize() uint {
	switch runtime.GOOS {
	case "darwin":
		return platform.GetPortRangeSizeOnMac()
	case "windows":
		return platform.GetPortRangeSizeOnWindows()
	default:
		return platform.GetPortRangeSizeOnLinux()
	}
}

func getTotalFileDescriptors() uint {
	switch runtime.GOOS {
	case "darwin":
		return platform.GetTotalFileDescriptorsOnMac()
	case "windows":
		return platform.GetTotalFileDescriptorsOnWindows()
	default:
		return platform.GetTotalFileDescriptorsOnLinux()
	}
}
