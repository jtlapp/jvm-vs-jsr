package platform

// These methods are untested.

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetPortsInUseCountsOnWindows() (timeWaitCount, establishedCount uint) {
	cmd := exec.Command("netstat", "-n", "-p", "TCP")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(out), "\n")

	// Skip header lines (usually 4 on Windows)
	for i := 4; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				state := fields[3]
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

func GetPortRangeSizeOnWindows() uint {
	cmd := exec.Command("netsh", "int", "ipv4", "show", "dynamicport", "tcp")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(out), "\n")
	var numPorts int

	for _, line := range lines {
		if strings.Contains(line, "Number of Ports") {
			fields := strings.Fields(line)
			numPorts, _ = strconv.Atoi(fields[len(fields)-1])
			break
		}
	}
	return uint(numPorts)
}
