package cmd

import (
	"flag"
	"fmt"
	"runtime"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/cli"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

var LoopDeterminingRates = cli.NewCommand(
	"loop",
	"-scenario=<scenario> [-times <iterations>] [<trial-options>]",
	"Loops repeatedly performing tests to find the highest constant/stable rate. "+
		"The resulting rates are guaranteed to be error-free for the specified "+
		"duration. Provide a rate guess to hasten convergence on the stable rate.",
	addLoopOptions,
	func(commandConfig config.CommandConfig) error {

		runCount := *commandConfig.LoopCount
		resetRandomSeed := *commandConfig.ResetRandomSeed

		runStats, err := performRuns(&commandConfig, runCount, resetRandomSeed)
		if err != nil {
			return err
		}

		util.Log()
		runStats.Print()
		return nil
	})

var DetermineRate = cli.NewCommand(
	"run",
	"-scenario=<scenario> [<trial-options>]",
	"Finds the highest constant/stable rate. The resulting rate is guaranteed "+
		"to be error-free for the specified duration. Provide a rate guess to hasten "+
		"convergence on the stable rate.",
	addTrialOptions,
	func(commandConfig config.CommandConfig) error {

		runCount := 1
		resetRandomSeed := false

		_, err := performRuns(&commandConfig, runCount, resetRandomSeed)
		return err
	})

var TryRate = cli.NewCommand(
	"try",
	"-scenario=<scenario> [<trial-options>]",
	"Tries issuing requests at the given rate for the specified duration.",
	addTrialOptions,
	func(commandConfig config.CommandConfig) error {

		resultsDB := database.NewResultsDatabase()
		defer resultsDB.Close()

		benchmarkRunner, err := createBenchmarkRunner(&commandConfig, resultsDB)
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

var ShowStatus = cli.NewCommand(
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
	commandConfig *config.CommandConfig,
	runCount int,
	resetRandomSeed bool,
) (*stats.RunStats, error) {

	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	benchmarkRunner, err := createBenchmarkRunner(commandConfig, resultsDB)
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

	config.LongSleepMillis = flagSet.Int("longSleepMillis", 400,
		"Duration of a long sleep request (in sleep scenarios)")

	config.ShortSleepMillis = flagSet.Int("shortSleepMillis", 100,
		"Duration of a short sleep request (in sleep scenarios)")

	config.PercentLongRequests = flagSet.Int("percentLongRequests", 20,
		"Percentage of requests that are long sleep requests (in sleep scenarios)")

	config.ConfigFile = cli.AllowConfigFile(flagSet)
}

func createBenchmarkRunner(
	commandConfig *config.CommandConfig,
	resultsDB *database.ResultsDB,
) (*runner.BenchmarkRunner, error) {

	platformConfig, err := config.GetPlatformConfig()
	if err != nil {
		return nil, err
	}

	scenario, err := scenarios.GetScenario(*commandConfig.ScenarioName)
	if err != nil {
		return nil, err
	}

	return runner.NewBenchmarkRunner(
		*platformConfig,
		*commandConfig,
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
