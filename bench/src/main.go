package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	// Import test suites

	vegeta "github.com/tsenart/vegeta/lib"
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

type TestSuite interface {
	Name() string
	PerformSetup() error
	GetTargeter(baseUrl string) vegeta.Targeter
}

var testSuitesSlice = []TestSuite{
	&sleep.Suite{},
	&taggedints.Suite{},
}

func getTestSuite(name string) (TestSuite, bool) {
	for _, suite := range testSuitesSlice {
		if suite.Name() == name {
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

	if config.mode == "setup" {
		err := suite.PerformSetup()
		if err != nil {
			fail("Setup failed: %v", err)
		}
	} else if config.mode == "test" {
		runBenchmark(config, suite)
	} else {
		fail("Invalid argument '%s'. Must be 'setup' or 'test'.", config.mode)
	}
}

func parseArgs() Config {
	if len(os.Args) < 4 {
		fail("Usage: %s <test-suite-name> setup|test <rate> <duration>", os.Args[0])
	}

	suiteName := os.Args[1]
	mode := os.Args[2]

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	rate := flagSet.Int("rate", 10, "Requests per second")
	duration := flagSet.Int("duration", 5, "Duration of the benchmark in seconds")
	flagSet.Parse(os.Args[3:])

	baseUrl := os.Getenv(baseUrlEnvVar)
	if baseUrl == "" {
		fail("%s environment variable is required", baseUrlEnvVar)
	}

	return Config{baseUrl, suiteName, mode, *rate, *duration}
}

func runBenchmark(config Config, suite TestSuite) {

	targeter := suite.GetTargeter(config.baseUrl)

	attacker := vegeta.NewAttacker()
	rateLimiter := vegeta.Rate{Freq: config.rate, Per: time.Second}
	duration := time.Duration(config.durationSeconds) * time.Second

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rateLimiter, duration,
		"Benchmark sleep API") {
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
