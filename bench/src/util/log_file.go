package util

import (
	"fmt"
	"os"
)

const (
	logfileName = "benchmark.log"
)

func Logf(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	Log(line)
}

func Log(values ...interface{}) {
	// Write the values as a line to stdout.

	line := fmt.Sprint(values...)
	fmt.Println(line)

	// Write the values as a line to stdout.

	file, err := os.OpenFile(logfileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file")
	}
	defer file.Close()

	if _, err := file.WriteString(line + "\n"); err != nil {
		panic("failed to write to log file")
	}
}
