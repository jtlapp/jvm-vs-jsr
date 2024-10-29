package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"jvm-vs-jsr.jtlapp.com/benchmark/backend"
	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/orderitems"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/sleep"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/taggedints"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	version = "0.1.0"
	baseAppUrlEnvVar = "BASE_APP_URL"
)

var scenariosSlice = []runner.Scenario{
	&sleep.Scenario{},
	&taggedints.Scenario{},
	&orderitems.Scenario{},
}

func getScenario(name string) (runner.Scenario, bool) {
	for _, scenario := range scenariosSlice {
		if scenario.GetName() == name {
			return scenario, true
		}
	}
	return nil, false
}

func main() {
	if len(os.Args) == 1 {
		showUsage()
		os.Exit(0)
	}
	command := os.Args[1]

	backendDB := backend.NewBackendDatabase()
	defer backendDB.ClosePool()

	switch command {
	case "setup":
		scenario := parseScenario(backendDB)
		if err := scenario.SetUpTestTables(); err != nil {
			fail("Failed to set up DB: %v", err)
		}
		if err := scenario.SetSharedQueries(); err != nil {
			fail("Failed to set queries: %v", err)
		}
	case "set-queries":
		scenario := parseScenario(backendDB)
		if err := scenario.SetSharedQueries(); err != nil {
			fail("Failed to set queries: %v", err)
		}
	case "run":
		util.LogCommand()
		scenario := parseScenario(backendDB)
		benchmarkConfig := parseBenchmarkArgs(scenario.GetName())

		benchmarkStats := runner.NewBenchmarkRunner(benchmarkConfig, scenario).DetermineRate()
		util.Log("")
		benchmarkStats.Print()
		util.Log("CPUs used: %d", benchmarkConfig.CPUCount)
	case "test":
		util.LogCommand()
		scenario := parseScenario(backendDB)
		benchmarkConfig := parseBenchmarkArgs(scenario.GetName())

		metrics := runner.NewBenchmarkRunner(benchmarkConfig, scenario).TestRate(
			benchmarkConfig.InitialRate, benchmarkConfig.DurationSeconds)
		util.Log("")
		runner.PrintMetrics(metrics)
		util.Log("CPUs used: %d", benchmarkConfig.CPUCount)
	case "status":
		timeWaitPercent, establishedPercent := util.GetPortsInUsePercents()
		fmt.Printf("  active ports: %d%%, waiting ports: %d%%, FDs in use: %d%%\n\n",
			establishedPercent, timeWaitPercent, util.GetFDsInUsePercent())
	default:
		fail("Invalid argument command '%s'", command)
	}
}

func parseScenario(backendDB *backend.BackendDB) runner.Scenario {
	if len(os.Args) < 3 {
		failWithUsage("Scenario name is required")
	}
	scenarioName := os.Args[2]

	scenario, valid := getScenario(scenarioName)
	if !valid {
		fail("Unknown test scenario: %s", scenarioName)
	}
	if err := scenario.Init(backendDB); err != nil {
		fail("Initialization failed: %v", err)
	}
	return scenario
}

func parseBenchmarkArgs(scenarioName string) runner.BenchmarkConfig {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cpuCount := flagSet.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	maxConnections := flagSet.Int("maxconns", 0, "Maximum number of connections to use")
	rate := flagSet.Int("rate", 10, "Requests per second")
	duration := flagSet.Int("duration", 5, "Duration of the benchmark in seconds")
	timeout := flagSet.Int("timeout", 10, "Request response timeout in seconds")
	minWait := flagSet.Int("minwait", 0, "Minimum wait time between tests in seconds")
	if len(os.Args) > 3 {
		flagSet.Parse(os.Args[3:])
	}

	baseAppUrl := os.Getenv(baseAppUrlEnvVar)
	if baseAppUrl == "" {
		fail("%s environment variable is required", baseAppUrlEnvVar)
	}

	appInfo, err := util.GetAppInfo(baseAppUrl)
	if err != nil {
		fail("Failed to get app info: %v", err)
	}

	return runner.BenchmarkConfig{
		ClientVersion:         version,
		BaseAppUrl:            baseAppUrl,
		AppName:               appInfo.AppName,
		AppVersion:            appInfo.AppVersion,
		AppConfig:             appInfo.AppConfig,
		ScenarioName:          scenarioName,
		CPUCount:              *cpuCount,
		MaxConnections:        *maxConnections,
		InitialRate:           *rate,
		DurationSeconds:       *duration,
		RequestTimeoutSeconds: *timeout,
		MinWaitSeconds:        *minWait,
	}
}

func fail(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}

func failWithUsage(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	showUsage()
	os.Exit(1)
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
	for _, scenario := range scenariosSlice {
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
