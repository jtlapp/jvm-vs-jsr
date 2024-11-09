package command

import (
	"fmt"

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

	testResults, err := benchmarkRunner.DetermineRate()
	if err != nil {
		return err
	}
	util.Log()
	testResults.Print()
	return nil
}

func TestRate(clientConfig config.ClientConfig, argsParser *ArgsParser) error {
	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	benchmarkRunner, err := createBenchmarkRunner(clientConfig, argsParser, resultsDB)
	if err != nil {
		return err
	}

	metrics, err := benchmarkRunner.TestRate()
	if err != nil {
		return err
	}
	util.Log()
	runner.PrintMetrics(*metrics)
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
