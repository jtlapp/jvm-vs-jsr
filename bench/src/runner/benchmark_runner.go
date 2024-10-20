package runner

import (
	"fmt"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	warmupSeconds = 5
	resetSeconds  = 5
)

type BenchmarkRunner struct {
	config         BenchmarkConfig
	scenario       Scenario
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

func (br *BenchmarkRunner) DetermineRate() BenchmarkStats {

	// Warm up the application, in case it does JIT.

	fmt.Print("Warming up...")
	br.TestRate(br.config.InitialRate/5, warmupSeconds)

	// Find the highest rate that the system can handle without errors.

	rateUpperBound := br.config.InitialRate
	rateLowerBound := 0
	currentRate := rateUpperBound
	nextRate := -1

	for currentRate != 0 && nextRate != currentRate {

		time.Sleep(resetSeconds * time.Second)
		fmt.Printf("Testing %d requests/sec...\n", currentRate)
		br.currentMetrics = br.TestRate(currentRate, br.config.DurationSeconds)

		if br.currentMetrics.Success < 1 {
			rateUpperBound = currentRate
			if currentRate == rateLowerBound {
				rateLowerBound /= 2
				nextRate = rateLowerBound
			} else {
				nextRate = (rateLowerBound + rateUpperBound) / 2
			}
		} else {
			rateLowerBound = currentRate
			if currentRate == rateUpperBound {
				rateUpperBound *= 2
				nextRate = rateUpperBound
			} else {
				nextRate = (rateLowerBound + rateUpperBound) / 2
			}
		}
	}

	return BenchmarkStats{
		SteadyStateRate: currentRate,
		Metrics:         br.currentMetrics,
	}
}

func (br *BenchmarkRunner) TestRate(rate int, durationSeconds int) vegeta.Metrics {

	targetProvider := br.scenario.GetTargetProvider(br.config.BaseURL)

	attacker := vegeta.NewAttacker(vegeta.Workers(uint64(br.config.CPUCount)))
	rateLimiter := vegeta.Rate{Freq: rate, Per: time.Second}
	duration := time.Duration(durationSeconds) * time.Second

	var metrics vegeta.Metrics
	var trialName = fmt.Sprintf("Scenario %s, rate %d, duration %d", br.config.ScenarioName,
		rate, durationSeconds)
	for res := range attacker.Attack(targetProvider, rateLimiter, duration, trialName) {
		br.logger.Log(res.Code, string(res.Body))
		metrics.Add(res)
	}
	metrics.Close()
	return metrics
}
