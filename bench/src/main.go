package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/command"
	cmd "jvm-vs-jsr.jtlapp.com/benchmark/command"
	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	version          = "0.1.0"
	baseAppUrlEnvVar = "BASE_APP_URL"
	helpOption       = "-help"
)

func main() {

	// Extract environment variables.

	baseAppUrl := os.Getenv(baseAppUrlEnvVar)
	if baseAppUrl == "" {
		err := fmt.Errorf("%s environment variable is required", baseAppUrlEnvVar)
		if err != nil {
			fail(err)
		}
	}

	// Extract the command or show help.

	if len(os.Args) == 1 || os.Args[1] == helpOption {
		showUsage()
		os.Exit(0)
	}
	commandName := os.Args[1]
	command, err := cmd.Find(commandName)
	if err != nil {
		fail(err)
	}

	// Show command-specific help if requested.

	index := slices.IndexFunc(os.Args, func(arg string) bool {
		return strings.HasSuffix(arg, helpOption)
	})
	if index != -1 {
		command.PrintUsageWithOptions()
		os.Exit(0)
	}

	// Log the command line.

	commandLine := os.Args[0]
	for _, arg := range os.Args[1:] {
		commandLine += " " + arg
	}

	util.LogOnly("\n========================================")
	util.LogfOnly("\n%s %s", time.Now().Format("2006-01-02 15:04:05"), commandLine)

	// Execute the command.

	clientConfig := config.ClientConfig{ClientVersion: version, BaseAppUrl: baseAppUrl}
	commandConfig, err := command.ParseArgs()
	if err != nil {
		fail(err)
	}
	err = command.Execute(clientConfig, *commandConfig)
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	if err != nil {
		msg := fmt.Sprintf("Error: %v", err)
		util.Log(msg)
		if usage.IsUsageError(err) {
			showUsage()
		}
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println()
	fmt.Printf("Benchmark tool for testing the performance of a web application (v%s).\n", version)
	fmt.Println()
	fmt.Println("Commands (use --help to see arguments):")
	for _, cmd := range command.Commands {
		cmd.PrintUsage()
	}
	fmt.Println()
	fmt.Println("Available scenarios:")
	for _, scenario := range scenarios.GetScenarios() {
		fmt.Printf("    %s\n", scenario.GetName())
	}
	fmt.Println()
}
