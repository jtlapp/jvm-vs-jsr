package main

import (
	"fmt"
	"os"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/command"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	version          = "0.1.0"
	baseAppUrlEnvVar = "BASE_APP_URL"
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

	// Extract the command.

	if len(os.Args) == 1 {
		showUsage()
		os.Exit(0)
	}
	commandName := os.Args[1]

	// Log the command line.

	commandLine := os.Args[0]
	for _, arg := range os.Args[1:] {
		commandLine += " " + arg
	}

	util.Log()
	util.Log("========================================")
	util.Log()
	util.FLog("%s %s", time.Now().Format("2006-01-02 15:04:05"), commandLine)

	// Execute the command.

	clientConfig := config.ClientConfig{ClientVersion: version, BaseAppUrl: baseAppUrl}
	argsParser := command.NewArgsParser(clientConfig)
	var err error

	switch commandName {
	case "setup-backend":
		err = command.SetupBackendDB(argsParser)
	case "assign-queries":
		err = command.AssignQueries(argsParser)
	case "run":
		err = command.DetermineRate(clientConfig, argsParser)
	case "test":
		err = command.TestRate(clientConfig, argsParser)
	case "status":
		err = command.ShowStatus()
	default:
		err = fmt.Errorf("invalid argument command '%s'", commandName)
	}

	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	if err != nil {
		msg := fmt.Sprintf("Error: %v", err)
		util.Log(msg)
		fmt.Fprintf(os.Stderr, "%s\n", msg)
		if command.IsUsageError(err) {
			showUsage()
		}
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Printf("\nBenchmark tool for testing the performance of a web application (v%s).", version)

	fmt.Println("\nCommands:")
	fmt.Println("    setup-backend <scenario>")
	fmt.Println("        Creates database tables and queries required for the test scenario.")
	fmt.Println("    assign-queries <scenario>")
	fmt.Println("        Sets only the queries required for the test scenario")
	fmt.Println("    run <scenario> [<attack-options>]")
	fmt.Println("        Finds the highest constant/stable rate. The resulting rate is guaranteed")
	fmt.Println("      	 to be error-free for the specified duration. Provide a rate guess to hasten")
	fmt.Println("        convergence on the stable rate.")
	fmt.Println("    test <scenario> [<attack-options>]")
	fmt.Println("        Tests issuing requests at the given rate for the specified duration.")
	fmt.Println("    status")
	fmt.Println("        Prints the active ports, waiting ports, and file descriptors in use.")

	fmt.Println("\nAvailable scenarios:")
	for _, scenario := range scenarios.GetScenarios() {
		fmt.Printf("    %s\n", scenario.GetName())
	}

	fmt.Println("\nAttack options:")
	fmt.Println("    -cpus <number-of-CPUs>")
	fmt.Println("        Number of CPUs (and workers) to use (default: all CPUs)")
	fmt.Println("    -maxconns <number-of-connections>")
	fmt.Println("        Maximum number of connections to use (default: unlimited)")
	fmt.Println("    -rate <requests-per-second>")
	fmt.Println("        Rate to test or initial rate guess (default: 10)")
	fmt.Println("    -duration <seconds>")
	fmt.Println("        Test duration or time over which rate must be error-free (default: 5)")
	fmt.Println("    -timeout <seconds>")
	fmt.Println("        Request response timeout (default 10)")
	fmt.Println("    -minwait <seconds>")
	fmt.Println("        Minimum wait time between tests (default 0)")
	fmt.Println()
}
