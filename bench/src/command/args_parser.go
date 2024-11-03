package command

import (
	"flag"
	"os"
	"runtime"

	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type ArgsParser struct {
	clientInfo ClientInfo
}

func NewArgsParser(clientInfo ClientInfo) *ArgsParser {
	return &ArgsParser{clientInfo}
}

func (ap *ArgsParser) GetScenarioArg() (string, error) {
	if len(os.Args) < 3 {
		return "", NewUsageError("scenario name is required")
	}
	return os.Args[2], nil
}

func (ap *ArgsParser) GetBenchmarkArgs(scenarioName string) (*runner.BenchmarkConfig, error) {
	cpusPerNode := runtime.NumCPU()

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cpusToUse := flagSet.Int("cpus", cpusPerNode, "Number of CPUs to use")
	maxConnections := flagSet.Int("maxconns", 0, "Maximum number of connections to use")
	rate := flagSet.Int("rate", 10, "Requests per second")
	duration := flagSet.Int("duration", 5, "Duration of the benchmark in seconds")
	timeout := flagSet.Int("timeout", 10, "Request response timeout in seconds")
	minWait := flagSet.Int("minwait", 0, "Minimum wait time between tests in seconds")
	if len(os.Args) > 3 {
		flagSet.Parse(os.Args[3:])
	}

	appInfo, err := util.GetAppInfo(ap.clientInfo.BaseAppUrl)
	if err != nil {
		return nil, NewUsageError("Failed to get app info: %v", err)
	}

	return &runner.BenchmarkConfig{
		ClientVersion:         ap.clientInfo.ClientVersion,
		BaseAppUrl:            ap.clientInfo.BaseAppUrl,
		AppName:               appInfo.AppName,
		AppVersion:            appInfo.AppVersion,
		AppConfig:             appInfo.AppConfig,
		ScenarioName:          scenarioName,
		CPUsPerNode:           cpusPerNode,
		CPUsToUse:             *cpusToUse,
		WorkerCount:           *cpusToUse,
		MaxConnections:        *maxConnections,
		InitialRate:           *rate,
		DurationSeconds:       *duration,
		RequestTimeoutSeconds: *timeout,
		MinWaitSeconds:        *minWait,
	}, nil
}
