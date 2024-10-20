package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	// Import test scenarios

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
		benchmarkStats := runner.NewBenchmarkRunner(benchmarkConfig, scenario).DetermineRate(
			commandConfig.rate, commandConfig.durationSeconds,
		)
		fmt.Printf("CPUs used: %d\n", commandConfig.cpuCount)
		benchmarkStats.Print()
	case "test":
		benchmarkConfig := toBenchmarkConfig(commandConfig)
		benchmarkStats := runner.NewBenchmarkRunner(benchmarkConfig, scenario).TestRate(
			commandConfig.rate, commandConfig.durationSeconds,
		)
		fmt.Printf("CPUs used: %d\n", commandConfig.cpuCount)
		benchmarkStats.Print()
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
	if len(os.Args) > 3 {
		flagSet.Parse(os.Args[3:])
	}

	return commandConfig{
		scenarioName:    scenarioName,
		mode:            mode,
		cpuCount:        *cpuCount,
		rate:            *rate,
		durationSeconds: *duration,
	}
}

func toBenchmarkConfig(config commandConfig) runner.BenchmarkConfig {
	baseUrl := os.Getenv(baseUrlEnvVar)
	if baseUrl == "" {
		fail("%s environment variable is required", baseUrlEnvVar)
	}
	return runner.BenchmarkConfig{
		BaseURL:      baseUrl,
		ScenarioName: config.scenarioName,
		CPUCount:     config.cpuCount,
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
	fmt.Printf("\nUsage: %s <test-scenario-name> setup-all | set-queries | test\n", os.Args[0])
	fmt.Println("\nTest scenarios:")
	for _, scenario := range scenariosSlice {
		fmt.Printf("    %s\n", scenario.GetName())
	}
	fmt.Println("\n'run' finds the highest constant/stable rate. Options:")
	fmt.Println("    -cpus <number-of-CPUs>")
	fmt.Println("    -rate <requests-per-second> -- initial rate guess for hastening convergence")
	fmt.Println("    -duration <seconds> -- time over which rate must be error-free")
	fmt.Println()
	fmt.Println("\n'test' tests a provided rate. Options:")
	fmt.Println("    -cpus <number-of-CPUs>")
	fmt.Println("    -rate <requests-per-second> -- rate to test")
	fmt.Println("    -duration <seconds> -- duration of the test")
	fmt.Println()
}
