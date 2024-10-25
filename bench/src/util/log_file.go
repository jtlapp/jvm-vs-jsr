package util

import (
	"fmt"
	"os"
	"time"
)

const (
	logfileName = "benchmark.log"
)

func LogCommand() {
	commandLine := os.Args[0]
	for _, arg := range os.Args[1:] {
		commandLine += " " + arg
	}

	writeToFile("\n========================================\n")
	writeToFile(time.Now().Format("2006-01-02 15:04:05") + " " + commandLine)
	writeToFile("")
}

func Log(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	fmt.Println(line)
	writeToFile(line)
}

func writeToFile(line string) {
	file, err := os.OpenFile(logfileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file")
	}
	defer file.Close()

	if _, err := file.WriteString(line + "\n"); err != nil {
		panic("failed to write to log file")
	}
}
