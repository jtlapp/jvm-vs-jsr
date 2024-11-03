package main

import (
	"fmt"
	"os"

	"jvm-vs-jsr.jtlapp.com/benchmark/command"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	version          = "0.1.0"
	baseAppUrlEnvVar = "BASE_APP_URL"
)

func main() {
	if len(os.Args) == 1 {
		showUsage()
		os.Exit(0)
	}
	commandName := os.Args[1]
	util.LogCommand()
	var err error

	baseAppUrl := os.Getenv(baseAppUrlEnvVar)
	if baseAppUrl == "" {
		err = fmt.Errorf("%s environment variable is required", baseAppUrlEnvVar)
	}
	clientInfo := command.ClientInfo{ClientVersion: version, BaseAppUrl: baseAppUrl}
	argsParser := command.NewArgsParser(clientInfo)

	if err == nil {
		switch commandName {
		case "setup-backend":
			err = command.SetupBackendDB(argsParser)
		case "assign-queries":
			err = command.AssignQueries(argsParser)
		case "run":
			err = command.DetermineRate(argsParser)
		case "test":
			err = command.TestRate(argsParser)
		case "status":
			err = command.ShowStatus()
		default:
			err = fmt.Errorf("invalid argument command '%s'", commandName)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if command.IsUsageError(err) {
			showUsage()
		}
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Printf("\nBenchmark tool for testing the performance of a web application (v%s).", version)

	fmt.Println("\nCommands:")
	fmt.Println("    setup <scenario>")
	fmt.Println("        Creates database tables and queries required for the test scenario.")
	fmt.Println("    set-queries <scenario>")
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
