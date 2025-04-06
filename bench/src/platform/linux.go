package platform

import (
	"fmt"
	"os"
	"strings"
)

func GetPortsInUseCountsOnLinux() (timeWaitCount, establishedCount uint) {
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

func GetPortRangeSizeOnLinux() uint {
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
