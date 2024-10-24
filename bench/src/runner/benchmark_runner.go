package runner

import (
	"fmt"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	warmupRequestsPerSecond = 100
	warmupSeconds           = 5
)

type BenchmarkRunner struct {
	config         BenchmarkConfig
	scenario       Scenario
	successMetrics vegeta.Metrics
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

	fmt.Print("Warming up...\n")
	br.TestRate(warmupRequestsPerSecond, warmupSeconds)

	// Find the highest rate that the system can handle without errors.

	rateUpperBound := br.config.InitialRate
	rateLowerBound := 0
	currentRate := -1
	nextRate := rateUpperBound

	for currentRate != 0 && nextRate != currentRate {
		util.WaitForPortsToClear()
		currentRate = nextRate
		if currentRate == rateLowerBound {
			break
		}

		fmt.Printf("Testing %d requests/sec...\n", currentRate)
		metrics := br.TestRate(currentRate, br.config.DurationSeconds)
		printTestStatus(metrics)

		if metrics.Success < 1 {
			rateUpperBound = currentRate
			nextRate = (rateLowerBound + rateUpperBound) / 2
		} else {
			br.successMetrics = metrics
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
		Metrics:         br.successMetrics,
	}
}

func (br *BenchmarkRunner) TestRate(rate int, durationSeconds int) vegeta.Metrics {

	targetProvider := br.scenario.GetTargetProvider(br.config.BaseURL)

	attacker := vegeta.NewAttacker(
		vegeta.Workers(uint64(br.config.CPUCount)),
		vegeta.Connections(br.config.MaxConnections),
		vegeta.Timeout(time.Duration(br.config.RequestTimeoutSeconds)*time.Second),
		vegeta.KeepAlive(true),
	)
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

func printTestStatus(metrics vegeta.Metrics) {

	timeWaitPct, establishedPct := util.GetPortsInUsePercents()
	errorMessages := strings.Join(metrics.Errors, ", ")
	if errorMessages == "" {
		errorMessages = "(none)"
	}

	fmt.Printf(
		"  Success rate: %.1f%%, req/sec: %.1f, ports active: %d%%, ports waiting: %d%%, FDs: %d%%, errors: %s\n",
		metrics.Success*100,
		metrics.Rate,
		establishedPct,
		timeWaitPct,
		util.GetFDsInUsePercent(),
		errorMessages,
	)
}
