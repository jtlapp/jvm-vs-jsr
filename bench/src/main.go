package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	// Import test suites

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-js.jtlapp.com/benchmark/lib"
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
	rate            int
	durationSeconds int
}

var testSuitesSlice = []lib.TestSuite{
	&sleep.Suite{},
	&taggedints.Suite{},
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

	suite, valid := getTestSuite(config.suiteName)
	if !valid {
		fail("Unknown test suite: %s", config.suiteName)
	}
	if err := suite.Init(); err != nil {
		fail("Initialization failed: %v", err)
	}

	switch config.mode {
	case "setup-all":
		if err := suite.SetUpDatabase(); err != nil {
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
	if len(os.Args) < 3 {
		failWithUsage("Too few arguments")
	}

	suiteName := os.Args[1]
	mode := os.Args[2]

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	rate := flagSet.Int("rate", 10, "Requests per second")
	duration := flagSet.Int("duration", 5, "Duration of the benchmark in seconds")
	if len(os.Args) > 3 {
		flagSet.Parse(os.Args[3:])
	}

	baseUrl := os.Getenv(baseUrlEnvVar)
	if baseUrl == "" {
		fail("%s environment variable is required", baseUrlEnvVar)
	}

	return Config{baseUrl, suiteName, mode, *rate, *duration}
}

func runBenchmark(config Config, suite lib.TestSuite) {

	targetProvider := suite.GetTargetProvider(config.baseUrl)
	logger := lib.NewResponseLogger()

	attacker := vegeta.NewAttacker()
	rateLimiter := vegeta.Rate{Freq: config.rate, Per: time.Second}
	duration := time.Duration(config.durationSeconds) * time.Second

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targetProvider, rateLimiter, duration,
		"Benchmark sleep API") {
		logger.Log(res.Code, string(res.Body))
		metrics.Add(res)
	}

	metrics.Close()

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
	fmt.Printf("Usage: %s <test-suite-name> setup-all|set-queries|test\n", os.Args[0])
	fmt.Println("'test' options:")
	fmt.Println("  -rate <requests-per-second>")
	fmt.Println("  -duration <seconds>")
	os.Exit(1)
}
