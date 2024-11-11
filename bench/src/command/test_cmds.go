package command

import (
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

func DetermineRate(clientConfig config.ClientConfig, argsParser *ArgsParser) error {
	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	benchmarkRunner, err := createBenchmarkRunner(clientConfig, argsParser, resultsDB)
	if err != nil {
		return err
	}

	metrics, err := benchmarkRunner.DetermineRate()
	if err != nil {
		return err
	}
	util.Log()
	printMetrics(metrics)
	return nil
}

func TryRate(clientConfig config.ClientConfig, argsParser *ArgsParser) error {
	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	benchmarkRunner, err := createBenchmarkRunner(clientConfig, argsParser, resultsDB)
	if err != nil {
		return err
	}

	metrics, err := benchmarkRunner.TryRate()
	if err != nil {
		return err
	}
	util.Log()
	printMetrics(metrics)
	return nil
}

func ShowStatus() error {
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
}

func createBenchmarkRunner(clientConfig config.ClientConfig, argsParser *ArgsParser, resultsDB *database.ResultsDB) (*runner.BenchmarkRunner, error) {
	platformConfig, err := config.GetPlatformConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	scenarioName, err := argsParser.GetScenarioName()
	if err != nil {
		return nil, err
	}

	scenario, err := scenarios.GetScenario(scenarioName)
	if err != nil {
		return nil, err
	}

	testConfig, err := argsParser.GetTestConfig()
	if err != nil {
		return nil, err
	}

	return runner.NewBenchmarkRunner(*platformConfig, *testConfig, &scenario, resultsDB)
}

func printMetrics(metrics *vegeta.Metrics) {
	util.FLog("Steady state rate: %.1f", metrics.Rate)
	util.FLog("Throughput: %f requests/sec", metrics.Throughput)
	util.FLog("Requests: %d", metrics.Requests)
	util.FLog("Success Percentage: %.2f%%", metrics.Success*100)
	util.FLog("Average Latency: %s", metrics.Latencies.Mean)
	util.FLog("99th Percentile Latency: %s", metrics.Latencies.P99)
	util.FLog("Max Latency: %s", metrics.Latencies.Max)
	util.FLog("Status Codes: %v", metrics.StatusCodes)
}
