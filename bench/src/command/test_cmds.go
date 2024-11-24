package command

import (
	"flag"
	"fmt"
	"runtime"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

var LoopDeterminingRates = newCommand(
	"loop",
	"-scenario=<scenario> [-times <iterations>] [<trial-options>]",
	"Loops repeatedly performing tests to find the highest constant/stable rate. "+
		"The resulting rates are guaranteed to be error-free for the specified "+
		"duration. Provide a rate guess to hasten convergence on the stable rate.",
	addLoopOptions,
	func(commandConfig config.CommandConfig) error {

		testConfig, err := getTestConfig(commandConfig)
		if err != nil {
			return err
		}
		runCount := *commandConfig.LoopCount
		resetRandomSeed := *commandConfig.ResetRandomSeed

		runStats, err := performRuns(*testConfig, &commandConfig, runCount, resetRandomSeed)
		if err != nil {
			return err
		}

		util.Log()
		runStats.Print()
		return nil
	})

var DetermineRate = newCommand(
	"run",
	"-scenario=<scenario> [<trial-options>]",
	"Finds the highest constant/stable rate. The resulting rate is guaranteed "+
		"to be error-free for the specified duration. Provide a rate guess to hasten "+
		"convergence on the stable rate.",
	addTrialOptions,
	func(commandConfig config.CommandConfig) error {

		testConfig, err := getTestConfig(commandConfig)
		if err != nil {
			return err
		}
		runCount := 1
		resetRandomSeed := false

		_, err = performRuns(*testConfig, &commandConfig, runCount, resetRandomSeed)
		return err
	})

var TryRate = newCommand(
	"try",
	"-scenario=<scenario> [<trial-options>]",
	"Tries issuing requests at the given rate for the specified duration.",
	addTrialOptions,
	func(commandConfig config.CommandConfig) error {

		testConfig, err := getTestConfig(commandConfig)
		if err != nil {
			return err
		}

		resultsDB := database.NewResultsDatabase()
		defer resultsDB.Close()

		benchmarkRunner, err := createBenchmarkRunner(*testConfig, &commandConfig, resultsDB)
		if err != nil {
			return err
		}

		metrics, err := benchmarkRunner.TryRate()
		if err != nil {
			return err
		}

		util.Log()
		printTrialMetrics(metrics)
		return nil
	})

var ShowStatus = newCommand(
	"status",
	"",
	"Prints the active ports, waiting ports, and file descriptors in use.",
	nil,
	func(commandConfig config.CommandConfig) error {
		resources := util.NewResourceStatus()
		establishedPortsPercent, timeWaitPortsPercent, fdsInUsePercent :=
			resources.GetPercentages()
		fmt.Printf("  active ports: %d of %d (%d%%)\n",
			resources.EstablishedPortsCount, resources.TotalAvailablePorts,
			uint(establishedPortsPercent+.5))
		fmt.Printf("  waiting ports: %d of %d (%d%%)\n",
			resources.TimeWaitPortsCount, resources.TotalAvailablePorts,
			uint(timeWaitPortsPercent+.5))
		fmt.Printf("  FDs in use: %d of %d (%d%%)\n",
			resources.FDsInUseCount, resources.TotalFileDescriptors,
			uint(fdsInUsePercent+.5))
		return nil
	})

func performRuns(
	testConfig config.TestConfig,
	commandConfig *config.CommandConfig,
	runCount int,
	resetRandomSeed bool,
) (*stats.RunStats, error) {

	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	benchmarkRunner, err := createBenchmarkRunner(testConfig, commandConfig, resultsDB)
	if err != nil {
		return nil, err
	}

	return benchmarkRunner.DetermineRate(runCount, resetRandomSeed)
}

func addLoopOptions(config *config.CommandConfig, flagSet *flag.FlagSet) {
	config.LoopCount = flagSet.Int("runCount", 8,
		"Number of times to run the rate determination benchmark")

	config.ResetRandomSeed = flagSet.Bool("resetSeedBetweenTests", false,
		"Reset the random seed for each run")

	addTrialOptions(config, flagSet)
}

func addTrialOptions(config *config.CommandConfig, flagSet *flag.FlagSet) {
	config.ScenarioName = flagSet.String("scenario", "",
		"Name of scenario to test (REQUIRED)")

	config.CPUsToUse = flagSet.Int("cpusToUse", runtime.NumCPU(),
		"Number of CPUs (and workers) to use")

	config.MaxConnections = flagSet.Int("maxConnections", 0,
		"Maximum number of connections to use (default 0, meaning unlimited)")

	config.InitialRequestsPerSecond = flagSet.Int("initialRate", 10,
		"Rate to test or initial rate guess in requests/second. Ignored when querying for statistics.")

	config.DurationSeconds = flagSet.Int("testDuration", 5,
		"Test duration in seconds. Time over which rate must be error-free.")

	config.RequestTimeoutSeconds = flagSet.Int("requestTimeout", 10,
		"Request response timeout in seconds.")

	config.MinWaitSeconds = flagSet.Int("minWaitBetweenTests", 0,
		"Minimum wait time between tests in seconds (default 0)")

	config.InitialRandomSeed = flagSet.Int("seed", 123456,
		"Random seed for randomizing requests (in supporting scenarios). When "+
			"querying for statistics, set to 0 to query across all random seeds.")
}

func getTestConfig(commandConfig config.CommandConfig) (*config.TestConfig, error) {
	scenarioName := *commandConfig.ScenarioName
	if scenarioName == "" {
		return nil, usage.NewUsageError("scenario name is required")
	}

	return &config.TestConfig{
		ScenarioName:             scenarioName,
		CPUsToUse:                *commandConfig.CPUsToUse,
		WorkerCount:              *commandConfig.CPUsToUse,
		MaxConnections:           *commandConfig.MaxConnections,
		InitialRequestsPerSecond: *commandConfig.InitialRequestsPerSecond,
		InitialRandomSeed:        *commandConfig.InitialRandomSeed,
		DurationSeconds:          *commandConfig.DurationSeconds,
		RequestTimeoutSeconds:    *commandConfig.RequestTimeoutSeconds,
		MinWaitSeconds:           *commandConfig.MinWaitSeconds,
	}, nil
}

func createBenchmarkRunner(
	testConfig config.TestConfig,
	commandConfig *config.CommandConfig,
	resultsDB *database.ResultsDB,
) (*runner.BenchmarkRunner, error) {

	platformConfig, err := config.GetPlatformConfig()
	if err != nil {
		return nil, err
	}

	scenario, err := scenarios.GetScenario(testConfig.ScenarioName)
	if err != nil {
		return nil, err
	}

	scenarioConfig := config.NewScenarioConfig(commandConfig)

	return runner.NewBenchmarkRunner(
		*platformConfig,
		testConfig,
		*scenarioConfig,
		&scenario,
		resultsDB)
}

func printTrialMetrics(metrics *vegeta.Metrics) {
	util.Logf("Steady state rate: %.1f req/sec", metrics.Rate)
	util.Logf("Successful completions: %f req/sec", metrics.Throughput)
	util.Logf("Requests: %d", metrics.Requests)
	util.Logf("Success Percentage: %.2f%%", metrics.Success*100)
	util.Logf("Average Latency: %s", metrics.Latencies.Mean)
	util.Logf("99th Percentile Latency: %s", metrics.Latencies.P99)
	util.Logf("Max Latency: %s", metrics.Latencies.Max)
	util.Logf("Status Codes: %v", metrics.StatusCodes)
}
