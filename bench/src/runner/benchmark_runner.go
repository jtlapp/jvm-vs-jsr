package runner

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type BenchmarkRunner struct {
	config         BenchmarkConfig
	scenario       Scenario
	rateUnderTest  int
	currentMetrics vegeta.Metrics
	logger         *ResponseLogger
}

func NewBenchmarkRunner(config BenchmarkConfig, scenario Scenario) *BenchmarkRunner {
	return &BenchmarkRunner{
		config:   config,
		scenario: scenario,
		logger:   NewResponseLogger(),
	}
}

func (br *BenchmarkRunner) DetermineRate(initialRate int, durationSeconds int) BenchmarkStats {

	return BenchmarkStats{
		SteadyStateRate: br.rateUnderTest,
		Metrics:         br.currentMetrics,
	}
}

func (br *BenchmarkRunner) TestRate(rate int, durationSeconds int) BenchmarkStats {

	targetProvider := br.scenario.GetTargetProvider(br.config.BaseURL)

	attacker := vegeta.NewAttacker(vegeta.Workers(uint64(br.config.CPUCount)))
	rateLimiter := vegeta.Rate{Freq: rate, Per: time.Second}
	duration := time.Duration(durationSeconds) * time.Second

	var metrics vegeta.Metrics
	var trialName = fmt.Sprintf("Scenario %s, rate %d, duration %d", br.config.ScenarioName, rate, durationSeconds)
	for res := range attacker.Attack(targetProvider, rateLimiter, duration, trialName) {
		br.logger.Log(res.Code, string(res.Body))
		metrics.Add(res)
	}
	metrics.Close()

	return BenchmarkStats{
		SteadyStateRate: rate,
		Metrics:         metrics,
	}
}
