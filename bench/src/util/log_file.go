package util

import (
	"fmt"
	"os"
)

const (
	logfileName = "benchmark.log"
)

func Log(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)

	// Write line to stdout

	fmt.Println(line)

	// Write line to log file

	file, err := os.OpenFile(logfileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file")
	}
	defer file.Close()

	if _, err := file.WriteString(line + "\n"); err != nil {
		panic("failed to write to log file")
	}
}
