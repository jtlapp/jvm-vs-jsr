package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	// Import test suites

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-js.jtlapp.com/benchmark/lib"
	"jvm-vs-js.jtlapp.com/benchmark/suites/orderitems"
	"jvm-vs-js.jtlapp.com/benchmark/suites/sleep"
	"jvm-vs-js.jtlapp.com/benchmark/suites/taggedints"
)

const (
	baseUrlEnvVar = "BASE_APP_URL"
)

type Config struct {
	baseUrl         string
	suiteName       string
	mode            string
	cpuCount        int
	rate            int
	durationSeconds int
}

var testSuitesSlice = []lib.TestSuite{
	&sleep.Suite{},
	&taggedints.Suite{},
	&orderitems.Suite{},
}

func getTestSuite(name string) (lib.TestSuite, bool) {
	for _, suite := range testSuitesSlice {
		if suite.GetName() == name {
			return suite, true
		}
	}
	return nil, false
}

func main() {
	config := parseArgs()

	backendDB := lib.NewBackendDatabase()
	defer backendDB.ClosePool()

	suite, valid := getTestSuite(config.suiteName)
	if !valid {
		fail("Unknown test suite: %s", config.suiteName)
	}
	if err := suite.Init(backendDB); err != nil {
		fail("Initialization failed: %v", err)
	}

	switch config.mode {
	case "setup-all":
		if err := suite.SetUpTestTables(); err != nil {
			fail("Failed to set up DB: %v", err)
		}
		if err := suite.SetSharedQueries(); err != nil {
			fail("Failed to set queries: %v", err)
		}
	case "set-queries":
		if err := suite.SetSharedQueries(); err != nil {
			fail("Failed to set queries: %v", err)
		}
	case "test":
		runBenchmark(config, suite)
	default:
		fail("Invalid argument '%s'. Must be 'setup' or 'test'.", config.mode)
	}
}

func parseArgs() Config {
	if len(os.Args) == 1 {
		showUsage()
		os.Exit(0)
	} else if len(os.Args) < 3 {
		failWithUsage("Too few arguments")
	}

	suiteName := os.Args[1]
	mode := os.Args[2]

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cpuCount := flagSet.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	rate := flagSet.Int("rate", 10, "Requests per second")
	duration := flagSet.Int("duration", 5, "Duration of the benchmark in seconds")
	if len(os.Args) > 3 {
		flagSet.Parse(os.Args[3:])
	}

	baseUrl := os.Getenv(baseUrlEnvVar)
	if baseUrl == "" {
		fail("%s environment variable is required", baseUrlEnvVar)
	}

	return Config{baseUrl, suiteName, mode, *cpuCount, *rate, *duration}
}

func runBenchmark(config Config, suite lib.TestSuite) {

	targetProvider := suite.GetTargetProvider(config.baseUrl)
	logger := lib.NewResponseLogger()

	attacker := vegeta.NewAttacker(vegeta.Workers(uint64(config.cpuCount)))
	rateLimiter := vegeta.Rate{Freq: config.rate, Per: time.Second}
	duration := time.Duration(config.durationSeconds) * time.Second

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targetProvider, rateLimiter, duration,
		"Benchmark sleep API") {
		logger.Log(res.Code, string(res.Body))
		metrics.Add(res)
	}

	metrics.Close()

	fmt.Printf("CPUs used: %d\n", config.cpuCount)
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("99th Percentile Latency: %s\n", metrics.Latencies.P99)
	fmt.Printf("Max Latency: %s\n", metrics.Latencies.Max)
	fmt.Printf("Status Codes: %v\n", metrics.StatusCodes)
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
	fmt.Printf("\nUsage: %s <test-suite-name> setup-all | set-queries | test\n", os.Args[0])
	fmt.Println("\nTest suites:")
	for _, suite := range testSuitesSlice {
		fmt.Printf("    %s\n", suite.GetName())
	}
	fmt.Println("\n'test' options:")
	fmt.Println("    -cpus <number-of-CPUs>")
	fmt.Println("    -rate <requests-per-second>")
	fmt.Println("    -duration <seconds>")
	fmt.Println()
}
