package platform

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetPortsInUseCountsOnMac() (timeWaitCount, establishedCount uint) {
	cmd := exec.Command("netstat", "-n", "-p", "tcp")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(out), "\n")

	// Skip header lines (usually 2 on macOS)
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				state := fields[5]
				switch state {
				case "TIME_WAIT":
					timeWaitCount++
				case "ESTABLISHED":
					establishedCount++
				}
			}
		}
	}
	return timeWaitCount, establishedCount
}

func GetPortRangeSizeOnMac() uint {
	cmd := exec.Command("sysctl", "-n", "net.inet.ip.portrange.first", "net.inet.ip.portrange.last")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var lowPort, highPort int
	_, err = fmt.Sscanf(string(out), "%d\n%d", &lowPort, &highPort)
	if err != nil {
		panic(err)
	}
	return uint(highPort - lowPort)
}
