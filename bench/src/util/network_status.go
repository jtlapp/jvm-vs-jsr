package util

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

var portRangeSize = getPortRangeSize()

func GetFDsInUsePercent() uint {
	inUseFDs, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		panic(err)
	}
	inUseFDCount := float64(len(inUseFDs))

	var rlimit syscall.Rlimit
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}

	// Round to the nearest percent.

	return uint((inUseFDCount/float64(rlimit.Cur))*100.0 + 0.5)
}

func GetPortsInUsePercents() (timeWaitPercent, establishedPercent int) {
	var timeWaitCount, establishedCount float64

	data, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		panic(err)
	}

	// Skip header line
	lines := strings.Split(string(data), "\n")[1:]

	// Column 4 contains the connection state in hex
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		state := fields[3]
		switch state {
		case "06":
			timeWaitCount++
		case "01":
			establishedCount++
		}
	}

	timeWaitPercent = int((timeWaitCount/float64(portRangeSize))*100.0 + 0.5)
	establishedPercent = int((establishedCount/float64(portRangeSize))*100.0 + 0.5)

	return timeWaitPercent, establishedPercent
}

func getPortRangeSize() uint16 {
	data, err := os.ReadFile("/proc/sys/net/ipv4/ip_local_port_range")
	if err != nil {
		panic(err)
	}

	var lowPort, highPort int
	_, err = fmt.Sscanf(string(data), "%d %d", &lowPort, &highPort)
	if err != nil {
		panic(err)
	}
	return uint16(highPort - lowPort)
}
