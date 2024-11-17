package command

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	loopCount            = "times"
	resetRandomSeed      = "reset"
	cpusOption           = "cpus"
	maxConnectionsOption = "maxconns"
	rateOption           = "rate"
	durationOption       = "duration"
	timeoutOption        = "timeout"
	minWaitOption        = "minwait"
	randomSeedOption     = "seed"
)

const (
	defaultLoopCount       = 8
	defaultResetRandomSeed = false
	defaultMaxConnections  = 0
	defaultRate            = 10
	defaultDuration        = 5
	defaultTimeout         = 10
	defaultMinWait         = 0
	defaultRandomSeed      = 123456
)

var LoopDeterminingRates = newCommand(
	"loop",
	"<scenario> [-times <iterations>] [<trial-options>]",
	"Loops repeatedly performing tests to find the highest constant/stable rate. "+
		"The resulting rates are guaranteed to be error-free for the specified "+
		"duration. Provide a rate guess to hasten convergence on the stable rate.",
	printLoopOptions,
	func(clientConfig config.ClientConfig) error {
		flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		runCount := flagSet.Int(loopCount, defaultLoopCount, "")
		resetRandomSeed := flagSet.Bool(resetRandomSeed, defaultResetRandomSeed, "")

		runStats, err := performRuns(clientConfig, flagSet, runCount, resetRandomSeed)
		if err != nil {
			return err
		}

		util.Log()
		runStats.Print()
		return nil
	})

var DetermineRate = newCommand(
	"run",
	"<scenario> [<trial-options>]",
	"Finds the highest constant/stable rate. The resulting rate is guaranteed "+
		"to be error-free for the specified duration. Provide a rate guess to hasten "+
		"convergence on the stable rate.",
	printTrialOptions,
	func(clientConfig config.ClientConfig) error {
		flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		runCount := 1
		resetRandomSeed := false
		_, err := performRuns(clientConfig, flagSet, &runCount, &resetRandomSeed)
		return err
	})

var TryRate = newCommand(
	"try",
	"<scenario> [<trial-options>]",
	"Tries issuing requests at the given rate for the specified duration.",
	printTrialOptions,
	func(clientConfig config.ClientConfig) error {
		flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		testConfig, err := getTestConfig(flagSet)
		if err != nil {
			return err
		}

		resultsDB := database.NewResultsDatabase()
		defer resultsDB.Close()

		benchmarkRunner, err := createBenchmarkRunner(clientConfig, testConfig, resultsDB)
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
	func(clientConfig config.ClientConfig) error {
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
	clientConfig config.ClientConfig,
	flagSet *flag.FlagSet,
	runCount *int,
	resetRandomSeed *bool,
) (*stats.RunStats, error) {

	testConfig, err := getTestConfig(flagSet)
	if err != nil {
		return nil, err
	}

	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	benchmarkRunner, err := createBenchmarkRunner(clientConfig, testConfig, resultsDB)
	if err != nil {
		return nil, err
	}

	return benchmarkRunner.DetermineRate(*runCount, *resetRandomSeed)
}

func getTestConfig(flagSet *flag.FlagSet) (*config.TestConfig, error) {
	scenarioName, err := usage.GetScenarioName()
	if err != nil {
		return nil, err
	}

	cpusToUse := flagSet.Int(cpusOption, runtime.NumCPU(), "")
	maxConnections := flagSet.Int(maxConnectionsOption, defaultMaxConnections, "")
	rate := flagSet.Int(rateOption, defaultRate, "")
	duration := flagSet.Int(durationOption, defaultDuration, "")
	timeout := flagSet.Int(timeoutOption, defaultTimeout, "")
	minWait := flagSet.Int(minWaitOption, defaultMinWait, "")
	randomSeed := flagSet.Int(randomSeedOption, defaultRandomSeed, "")

	if len(os.Args) > 3 {
		err := flagSet.Parse(os.Args[3:])
		if err != nil {
			return nil, usage.NewUsageError("%s", err.Error())
		}
	}

	return &config.TestConfig{
		ScenarioName:             scenarioName,
		CPUsToUse:                *cpusToUse,
		WorkerCount:              *cpusToUse,
		MaxConnections:           *maxConnections,
		InitialRequestsPerSecond: *rate,
		InitialRandomSeed:        *randomSeed,
		DurationSeconds:          *duration,
		RequestTimeoutSeconds:    *timeout,
		MinWaitSeconds:           *minWait,
	}, nil
}

func createBenchmarkRunner(
	clientConfig config.ClientConfig,
	testConfig *config.TestConfig,
	resultsDB *database.ResultsDB,
) (*runner.BenchmarkRunner, error) {

	platformConfig, err := config.GetPlatformConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	scenarioName, err := usage.GetScenarioName()
	if err != nil {
		return nil, err
	}

	scenario, err := scenarios.GetScenario(scenarioName)
	if err != nil {
		return nil, err
	}

	return runner.NewBenchmarkRunner(*platformConfig, *testConfig, &scenario, resultsDB)
}

func printLoopOptions() {
	usage.PrintOption(
		loopCount,
		"iterations",
		"Number of times to run the rate determination benchmark",
		strconv.Itoa(defaultLoopCount),
	)
	printTrialOptions()
}

func printTrialOptions() {
	usage.PrintOption(
		cpusOption,
		"number of CPUs",
		"Number of CPUs (and workers) to use",
		"all CPUs",
	)
	usage.PrintOption(
		maxConnectionsOption,
		"number of connections",
		"Maximum number of connections to use",
		"unlimited",
	)
	usage.PrintOption(
		rateOption,
		"requests per second",
		"Rate to test or initial rate guess. Ignored when querying for statistics.",
		strconv.Itoa(defaultRate),
	)
	usage.PrintOption(
		durationOption,
		"seconds",
		"Test duration or time over which rate must be error-free",
		strconv.Itoa(defaultDuration),
	)
	usage.PrintOption(
		timeoutOption,
		"seconds",
		"Request response timeout",
		strconv.Itoa(defaultTimeout),
	)
	usage.PrintOption(
		minWaitOption,
		"seconds",
		"Minimum wait time between tests",
		strconv.Itoa(defaultMinWait),
	)
	usage.PrintOption(
		randomSeedOption,
		"random seed",
		"Random seed for randomizing requests (for supporting scenarios). When "+
			"querying for statistics, set to 0 to query across all random seeds.",
		strconv.Itoa(defaultRandomSeed),
	)
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
