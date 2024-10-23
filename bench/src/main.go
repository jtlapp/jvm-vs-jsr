package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/orderitems"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/sleep"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/taggedints"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	baseUrlEnvVar = "BASE_APP_URL"
)

type commandConfig struct {
	scenarioName    string
	mode            string
	cpuCount        int
	rate            int
	durationSeconds int
	timeoutSeconds  int
}

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
	commandConfig := parseArgs()

	backendDB := util.NewBackendDatabase()
	defer backendDB.ClosePool()

	scenario, valid := getScenario(commandConfig.scenarioName)
	if !valid {
		fail("Unknown test scenario: %s", commandConfig.scenarioName)
	}
	if err := scenario.Init(backendDB); err != nil {
		fail("Initialization failed: %v", err)
	}

	switch commandConfig.mode {
	case "setup-all":
		if err := scenario.SetUpTestTables(); err != nil {
			fail("Failed to set up DB: %v", err)
		}
		if err := scenario.SetSharedQueries(); err != nil {
			fail("Failed to set queries: %v", err)
		}
	case "set-queries":
		if err := scenario.SetSharedQueries(); err != nil {
			fail("Failed to set queries: %v", err)
		}
	case "run":
		benchmarkConfig := toBenchmarkConfig(commandConfig)
		benchmarkStats := runner.NewBenchmarkRunner(benchmarkConfig, scenario).DetermineRate()
		fmt.Println()
		benchmarkStats.Print()
		fmt.Printf("CPUs used: %d\n", commandConfig.cpuCount)
	case "test":
		benchmarkConfig := toBenchmarkConfig(commandConfig)
		metrics := runner.NewBenchmarkRunner(benchmarkConfig, scenario).TestRate(
			commandConfig.rate, commandConfig.durationSeconds)
		fmt.Println()
		runner.PrintMetrics(metrics)
		fmt.Printf("CPUs used: %d\n", commandConfig.cpuCount)
	default:
		fail("Invalid argument '%s'. Must be 'setup' or 'test'.", commandConfig.mode)
	}
}

func parseArgs() commandConfig {
	if len(os.Args) == 1 {
		showUsage()
		os.Exit(0)
	} else if len(os.Args) < 3 {
		failWithUsage("Too few arguments")
	}

	scenarioName := os.Args[1]
	mode := os.Args[2]

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cpuCount := flagSet.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	rate := flagSet.Int("rate", 10, "Requests per second")
	duration := flagSet.Int("duration", 5, "Duration of the benchmark in seconds")
	timeout := flagSet.Int("timeout", 10, "Request response timeout in seconds")
	if len(os.Args) > 3 {
		flagSet.Parse(os.Args[3:])
	}

	return commandConfig{
		scenarioName:    scenarioName,
		mode:            mode,
		cpuCount:        *cpuCount,
		rate:            *rate,
		durationSeconds: *duration,
		timeoutSeconds:  *timeout,
	}
}

func toBenchmarkConfig(config commandConfig) runner.BenchmarkConfig {
	baseUrl := os.Getenv(baseUrlEnvVar)
	if baseUrl == "" {
		fail("%s environment variable is required", baseUrlEnvVar)
	}
	return runner.BenchmarkConfig{
		BaseURL:               baseUrl,
		ScenarioName:          config.scenarioName,
		InitialRate:           config.rate,
		DurationSeconds:       config.durationSeconds,
		CPUCount:              config.cpuCount,
		RequestTimeoutSeconds: config.timeoutSeconds,
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
	fmt.Printf("\nUsage: %s <test-scenario-name> setup-all | set-queries | run | test\n", os.Args[0])
	fmt.Println("\nTest scenarios:")
	for _, scenario := range scenariosSlice {
		fmt.Printf("    %s\n", scenario.GetName())
	}

	fmt.Println("\nCommands:")
	fmt.Println("    setup-all -- Creates database tables and queries required for the test scenario.")
	fmt.Println("    set-queries -- Sets only the queries required for the test scenario")
	fmt.Println("    run -- Finds the highest constant/stable rate. The resulting rate is guaranteed")
	fmt.Println("      to be error-free for the specified duration. Provide a rate guess to hasten")
	fmt.Println("      convergence on the stable rate.")
	fmt.Println("    test -- Tests issuing requests at the given rate for the specified duration.")

	fmt.Println("\nOptions:")

	fmt.Println("    -cpus <number-of-CPUs> (default: all CPUs)")
	fmt.Println("    -rate <requests-per-second> -- Rate to test or initial rate guess")
	fmt.Println("    -duration <seconds> -- Test duration or time over which rate must be error-free")
	fmt.Println("    -timeout <seconds> -- Request response timeout and half delay between tests (default 10)")
	fmt.Println()
}
