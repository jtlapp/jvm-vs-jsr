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

	benchmarkStats := benchmarkRunner.DetermineRate()
	util.Log("")
	benchmarkStats.Print()
	util.Log("CPUs used: %d", benchmarkRunner.GetTestConfig().CPUsToUse)
	return nil
}

func TestRate(clientConfig config.ClientConfig, argsParser *ArgsParser) error {
	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	benchmarkRunner, err := createBenchmarkRunner(clientConfig, argsParser, resultsDB)
	if err != nil {
		return err
	}

	metrics := benchmarkRunner.TestRate()
	util.Log("")
	runner.PrintMetrics(metrics)
	util.Log("CPUs used: %d", benchmarkRunner.GetTestConfig().CPUsToUse)
	return nil
}

func ShowStatus() error {
	resources := util.NewResourceStatus()
	establishedPortsPercent, timeWaitPortsPercent, fdsInUsePercent :=
		resources.GetPercentages()
	fmt.Printf("  active ports: %d%%, waiting ports: %d%%, FDs in use: %d%%\n\n",
		uint(establishedPortsPercent+.5),
		uint(timeWaitPortsPercent+.5),
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

	dbPool, err := resultsDB.GetPool()
	if err != nil {
		return nil, err
	}

	return runner.NewBenchmarkRunner(*platformConfig, *testConfig, &scenario, dbPool)
}
