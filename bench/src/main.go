package main

import (
	"flag"
	"fmt"
	"os"

	// Import test suites

	"jvm-vs-js.jtlapp.com/benchmark/suites/sleep"
	"jvm-vs-js.jtlapp.com/benchmark/suites/suite2"
)

const (
	baseUrlEnvVar = "BASE_URL"
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
	Setup()
	Test()
}

var testSuitesSlice = []TestSuite{
	&sleep.Suite{},
	&suite2.Suite{},
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
		suite.Setup()
	} else if config.mode == "test" {
		suite.Test()
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

	rate := flag.Int("rate", 10, "Requests per second")
	duration := flag.Int("duration", 5, "Duration of the benchmark in seconds")
	flag.Parse()

	baseUrl := os.Getenv(baseUrlEnvVar)
	if baseUrl == "" {
		fail("%s environment variable is required", baseUrlEnvVar)
	}

	return Config{baseUrl, suiteName, mode, *rate, *duration}
}

// func runBenchmark(config Config, suite TestSuite) {
// 	url := fmt.Sprintf("%s/api/sleep/%d", config.baseUrl, sleepDuration)

// 	rateLimiter := vegeta.Rate{Freq: config.rate, Per: time.Second}

// 	targeter := vegeta.NewStaticTargeter(vegeta.Target{
// 		Method: "GET",
// 		URL:    url,
// 		Body:   bytes.NewBuffer(nil).Bytes(),
// 	})

// 	attacker := vegeta.NewAttacker()

// 	durationInSeconds := time.Duration(config.durationSeconds) * time.Second
// 	var metrics vegeta.Metrics
// 	for res := range attacker.Attack(targeter, rateLimiter, durationInSeconds, "Benchmark sleep API") {
// 		metrics.Add(res)
// 	}

// 	metrics.Close()
// }

func fail(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}
