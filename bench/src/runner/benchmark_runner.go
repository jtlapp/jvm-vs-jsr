package runner

import (
	"fmt"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	warmupRequestsPerSecond = 100
	warmupSeconds           = 5
)

type BenchmarkRunner struct {
	platformConfig config.PlatformConfig
	testConfig     config.TestConfig
	scenario       *scenarios.Scenario
	resultsDB      *database.ResultsDB
	successMetrics vegeta.Metrics
	logger         *ResponseLogger
}

func NewBenchmarkRunner(platformConfig config.PlatformConfig, testConfig config.TestConfig, scenario *scenarios.Scenario, resultsDB *database.ResultsDB) (*BenchmarkRunner, error) {
	return &BenchmarkRunner{
		platformConfig: platformConfig,
		testConfig:     testConfig,
		scenario:       scenario,
		resultsDB:      resultsDB,
		logger:         NewResponseLogger(),
	}, nil
}

func (br *BenchmarkRunner) GetTestConfig() config.TestConfig {
	return br.testConfig
}

func (br *BenchmarkRunner) DetermineRate() (*database.TestResults, error) {

	// Warm up the application, in case it does JIT.

	util.Log()
	util.Log("Warmup run (ignored)...")
	_, err := br.performRateTrial(warmupRequestsPerSecond, warmupSeconds)
	if err != nil {
		return nil, fmt.Errorf("error performing warmup trial: %v", err)
	}

	// Find the highest rate that the system can handle without errors.

	rateUpperBound := br.testConfig.InitialRequestsPerSecond
	rateLowerBound := 0
	currentRate := -1
	nextRate := rateUpperBound
	testsPerformed := 0
	startTime := time.Now()

	for currentRate != 0 && nextRate != currentRate {
		br.waitBetweenTests()

		currentRate = nextRate
		if currentRate == rateLowerBound {
			break
		}

		util.Log()
		util.FLog("Testing %d requests/sec...", currentRate)
		metrics, err := br.performRateTrial(currentRate, br.testConfig.DurationSeconds)
		if err != nil {
			return nil, fmt.Errorf("error performing rate trial: %v", err)
		}
		printTestStatus(metrics)

		if metrics.Success < 1 {
			rateUpperBound = currentRate
			nextRate = (rateLowerBound + rateUpperBound) / 2
		} else {
			br.successMetrics = *metrics
			rateLowerBound = currentRate
			if currentRate == rateUpperBound {
				rateUpperBound *= 2
				nextRate = rateUpperBound
			} else {
				nextRate = (rateLowerBound + rateUpperBound) / 2
			}
		}
		testsPerformed++
	}

	testResults := &database.TestResults{
		TestsPerformed:       testsPerformed,
		TotalDurationSeconds: int(time.Since(startTime).Seconds()),
		Metrics:              br.successMetrics,
	}
	resources := util.NewResourceStatus()
	err = br.resultsDB.SaveResults("det", &br.platformConfig, &br.testConfig, testResults, &resources)
	if err != nil {
		return nil, fmt.Errorf("error saving rate determination results: %v", err)
	}

	return testResults, nil
}

func (br *BenchmarkRunner) TestRate() (*vegeta.Metrics, error) {
	return br.performRateTrial(br.testConfig.InitialRequestsPerSecond, br.testConfig.DurationSeconds)
}

func (br *BenchmarkRunner) performRateTrial(rate int, durationSeconds int) (*vegeta.Metrics, error) {

	targetProvider := (*br.scenario).GetTargetProvider(br.platformConfig.BaseAppUrl)

	attacker := vegeta.NewAttacker(
		vegeta.Workers(uint64(br.testConfig.WorkerCount)),
		vegeta.Connections(br.testConfig.MaxConnections),
		vegeta.Timeout(time.Duration(br.testConfig.RequestTimeoutSeconds)*time.Second),
		vegeta.KeepAlive(true),
	)
	rateLimiter := vegeta.Rate{Freq: rate, Per: time.Second}
	duration := time.Duration(durationSeconds) * time.Second

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targetProvider, rateLimiter, duration, "") {
		br.logger.Log(res.Code, string(res.Body))
		metrics.Add(res)
	}
	metrics.Close()

	testResults := &database.TestResults{
		TestsPerformed:       1,
		TotalDurationSeconds: int(metrics.Duration.Seconds()),
		Metrics:              metrics,
	}
	resources := util.NewResourceStatus()
	err := br.resultsDB.SaveResults("trial", &br.platformConfig, &br.testConfig, testResults, &resources)
	if err != nil {
		return nil, fmt.Errorf("error saving trial results: %v", err)
	}

	return &metrics, nil
}

func (br *BenchmarkRunner) waitBetweenTests() {
	start := time.Now()
	util.WaitForPortsToClear(br.platformConfig.InitialPortsInUse)
	elapsed := time.Since(start)
	minDuration := time.Duration(br.testConfig.MinWaitSeconds) * time.Second

	if remainingTime := minDuration - elapsed; remainingTime > 0 {
		time.Sleep(remainingTime)
	}
}

func printTestStatus(metrics *vegeta.Metrics) {
	resourceStatus := util.NewResourceStatus()
	establishedPortsPercent, timeWaitPortsPercent, fdsInUsePercent :=
		resourceStatus.GetPercentages()

	errorMessages := strings.Join(metrics.Errors, ", ")
	if errorMessages == "" {
		errorMessages = "(none)"
	}

	util.FLog(
		"  %.1f%% successful (%.1f req/s): issued %.1f req/s, %d%% ports active, %d%% ports waiting, %d%% FDs, errors: %s",
		metrics.Success*100,
		metrics.Throughput,
		metrics.Rate,
		uint(establishedPortsPercent+.5),
		uint(timeWaitPortsPercent+.5),
		uint(fdsInUsePercent+.5),
		errorMessages,
	)
}
