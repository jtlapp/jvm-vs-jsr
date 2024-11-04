package command

import (
	"flag"
	"os"
	"runtime"

	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

type ArgsParser struct {
	clientConfig config.ClientConfig
}

func NewArgsParser(clientConfig config.ClientConfig) *ArgsParser {
	return &ArgsParser{clientConfig}
}

func (ap *ArgsParser) GetScenarioName() (string, error) {
	if len(os.Args) < 3 {
		return "", NewUsageError("scenario name is required")
	}
	return os.Args[2], nil
}

func (ap *ArgsParser) GetTestConfig() (*config.TestConfig, error) {
	scenarioName, err := ap.GetScenarioName()
	if err != nil {
		return nil, err
	}

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cpusToUse := flagSet.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	maxConnections := flagSet.Int("maxconns", 0, "Maximum number of connections to use")
	rate := flagSet.Int("rate", 10, "Requests per second")
	duration := flagSet.Int("duration", 5, "Duration of the benchmark in seconds")
	timeout := flagSet.Int("timeout", 10, "Request response timeout in seconds")
	minWait := flagSet.Int("minwait", 0, "Minimum wait time between tests in seconds")
	if len(os.Args) > 3 {
		flagSet.Parse(os.Args[3:])
	}

	return &config.TestConfig{
		ScenarioName:             scenarioName,
		CPUsToUse:                *cpusToUse,
		WorkerCount:              *cpusToUse,
		MaxConnections:           *maxConnections,
		InitialRequestsPerSecond: *rate,
		DurationSeconds:          *duration,
		RequestTimeoutSeconds:    *timeout,
		MinWaitSeconds:           *minWait,
	}, nil
}
