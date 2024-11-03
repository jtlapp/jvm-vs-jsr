package command

import (
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

func DetermineRate(argsParser *ArgsParser) error {
	benchmarkConfig, scenario, err := createBenchmarkConfig(argsParser)
	if err != nil {
		return err
	}

	benchmarkStats := runner.NewBenchmarkRunner(*benchmarkConfig, scenario).DetermineRate()
	util.Log("")
	benchmarkStats.Print()
	util.Log("CPUs used: %d", benchmarkConfig.CPUsToUse)
	return nil
}

func TestRate(argsParser *ArgsParser) error {
	benchmarkConfig, scenario, err := createBenchmarkConfig(argsParser)
	if err != nil {
		return err
	}

	metrics := runner.NewBenchmarkRunner(*benchmarkConfig, scenario).TestRate()
	util.Log("")
	runner.PrintMetrics(metrics)
	util.Log("CPUs used: %d", benchmarkConfig.CPUsToUse)
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

func createBenchmarkConfig(argsParser *ArgsParser) (*runner.BenchmarkConfig, *scenarios.Scenario, error) {
	scenarioName, err := argsParser.GetScenarioArg()
	if err != nil {
		return nil, nil, err
	}
	scenario, err := scenarios.GetScenario(scenarioName)
	if err != nil {
		return nil, nil, err
	}
	benchmarkConfig, err := argsParser.GetBenchmarkArgs(scenarioName)
	if err != nil {
		return nil, nil, err
	}
	return benchmarkConfig, &scenario, nil
}
