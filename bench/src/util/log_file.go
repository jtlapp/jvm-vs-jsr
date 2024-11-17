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
	line := fmt.Sprint(values...)
	fmt.Println(line)
	writeToLogFile(line)
}

func LogfOnly(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	LogOnly(line)
}

func LogOnly(values ...interface{}) {
	line := fmt.Sprint(values...)
	writeToLogFile(line)
}

func writeToLogFile(line string) {
	file, err := os.OpenFile(logfileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file")
	}
	defer file.Close()

	if _, err := file.WriteString(line + "\n"); err != nil {
		panic("failed to write to log file")
	}
}
