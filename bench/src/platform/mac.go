package platform

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func GetFDsInUseCountOnMac() uint {
	cmd := exec.Command("lsof", "-p", fmt.Sprintf("%d", os.Getpid()))
	out, err := cmd.Output()
	if err != nil {
		// lsof might return error code 1 if process exists but has no open files
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return 0
		}
		panic(err)
	}

	// Count lines in output, subtract 1 for header
	lines := strings.Split(string(out), "\n")
	count := len(lines) - 1
	if count < 0 {
		count = 0
	}
	return uint(count)
}

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

func GetTotalFileDescriptorsOnMac() uint {
	var rlimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	return uint(rlimit.Cur)
}
