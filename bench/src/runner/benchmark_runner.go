package runner

import (
	"fmt"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
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

func (br *BenchmarkRunner) DetermineRate(runCount int, resetRandomSeed bool) (*stats.RunStats, error) {
	if err := br.performWarmupRun(); err != nil {
		return nil, err
	}

	randomSeed := int64(br.testConfig.InitialRandomSeed)
	startTime := time.Now()

	for i := 0; i < runCount; i++ {
		err := br.performRateDetermination(randomSeed)
		if err != nil {
			return nil, err
		}
		if !resetRandomSeed {
			randomSeed++
		}
	}

	trials, err := br.resultsDB.GetTrials(startTime, &br.platformConfig, &br.testConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting trials: %v", err)
	}

	runStats, err := stats.CalculateRunStats(trials)
	if err != nil {
		return nil, fmt.Errorf("error calculating run stats: %v", err)
	}
	return &runStats, nil
}

func (br *BenchmarkRunner) TryRate() (metrics *vegeta.Metrics, err error) {
	var runID int

	if err = br.waitForPortsToClear(); err != nil {
		return nil, err
	}
	runID, err = br.resultsDB.CreateRun(&br.platformConfig, &br.testConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating single-trial run: %v", err)
	}

	var trialID int
	trialID, metrics, err = br.performRateTrial(
		runID,
		br.testConfig.InitialRequestsPerSecond,
		int64(br.testConfig.InitialRandomSeed),
		br.testConfig.DurationSeconds,
	)
	if err != nil {
		return nil, fmt.Errorf("error performing rate trial: %v", err)
	}

	err = br.resultsDB.UpdateRun(runID, br.testConfig.DurationSeconds, trialID)
	if err != nil {
		return nil, fmt.Errorf("error updating single-trial run: %v", err)
	}

	return metrics, err
}

func (br *BenchmarkRunner) performWarmupRun() error {
	if err := br.waitForPortsToClear(); err != nil {
		return err
	}

	util.Log()
	util.Log("Warmup run (ignored)...")
	_, _, err := br.performRateTrial(
		0,
		warmupRequestsPerSecond,
		int64(br.testConfig.InitialRandomSeed),
		warmupSeconds,
	)
	if err != nil {
		return fmt.Errorf("error performing warmup trial: %v", err)
	}
	return nil
}

func (br *BenchmarkRunner) performRateDetermination(randomSeed int64) error {
	runID, err := br.resultsDB.CreateRun(&br.platformConfig, &br.testConfig)
	if err != nil {
		return fmt.Errorf("error creating run: %v", err)
	}

	// Find the highest rate of successfully completing requests that the system
	// can handle without any request errors or timeouts.

	requestRateUpperBound := br.testConfig.InitialRequestsPerSecond
	requestRateLowerBound := 0
	var lowerBoundMetrics vegeta.Metrics
	currentRequestRate := -1
	nextRequestRate := requestRateUpperBound
	var bestTrialID int
	var metrics *vegeta.Metrics
	startTime := time.Now()

	for currentRequestRate != 0 && nextRequestRate != currentRequestRate {
		br.waitBetweenTests()

		currentRequestRate = nextRequestRate
		if currentRequestRate == requestRateLowerBound {
			break
		}

		util.Log()
		util.Logf("Testing %d requests/sec...", currentRequestRate)
		var trialID int
		trialID, metrics, err = br.performRateTrial(
			runID, currentRequestRate, randomSeed, br.testConfig.DurationSeconds)
		if err != nil {
			return fmt.Errorf("error performing rate trial: %v", err)
		}
		printTestStatus(metrics)

		if metrics.Success < 1 || metrics.Throughput < lowerBoundMetrics.Throughput {
			requestRateUpperBound = currentRequestRate
			nextRequestRate = (requestRateLowerBound + requestRateUpperBound) / 2
		} else {
			bestTrialID = trialID
			lowerBoundMetrics = *metrics
			requestRateLowerBound = currentRequestRate
			if currentRequestRate == requestRateUpperBound {
				requestRateUpperBound *= 2
				nextRequestRate = requestRateUpperBound
			} else {
				nextRequestRate = (requestRateLowerBound + requestRateUpperBound) / 2
			}
		}
	}

	totalDurationSeconds := int(time.Since(startTime).Seconds())
	err = br.resultsDB.UpdateRun(runID, totalDurationSeconds, bestTrialID)
	if err != nil {
		return fmt.Errorf("error updating multi-trial run: %v", err)
	}

	util.Log()
	printRunMetrics(&lowerBoundMetrics)
	return nil
}

func (br *BenchmarkRunner) performRateTrial(
	runID,
	rate int,
	randomSeed int64,
	durationSeconds int,
) (int, *vegeta.Metrics, error) {

	targetProvider := (*br.scenario).GetTargetProvider(
		br.platformConfig.BaseAppUrl, randomSeed)

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
	attacker.Stop()
	metrics.Close()

	var trialID int
	if runID != 0 {
		resources := util.NewResourceStatus()
		trialInfo, err := stats.NewTrialInfo(&metrics, randomSeed)
		if err != nil {
			panic(fmt.Errorf("failed to create trial info: %v", err))
		}
		trialID, err = br.resultsDB.SaveTrial(runID, trialInfo, &resources)
		if err != nil {
			return 0, nil, fmt.Errorf("error saving trial: %v", err)
		}
	}
	return trialID, &metrics, nil
}

func (br *BenchmarkRunner) waitBetweenTests() {
	start := time.Now()
	util.WaitForPortsToTimeout()
	elapsed := time.Since(start)
	minDuration := time.Duration(br.testConfig.MinWaitSeconds) * time.Second

	if remainingTime := minDuration - elapsed; remainingTime > 0 {
		time.Sleep(remainingTime)
	}
}

func (br *BenchmarkRunner) waitForPortsToClear() error {
	portsAreReady, err := util.PortsAreReady(br.platformConfig.MaxReservedPorts)
	if err != nil {
		return err
	}
	if !portsAreReady {
		util.Log()
		util.Log("Waiting for ports to clear...")
		util.WaitForPortsToTimeout()
	}
	return nil
}

func printRunMetrics(metrics *vegeta.Metrics) {
	util.Logf("Best rate: %.1f req/sec (successfully completing %.1f req/sec)", metrics.Rate, metrics.Throughput)
	util.Log()
	util.Logf("  Requests: %d", metrics.Requests)
	util.Logf("  Average Latency: %s", metrics.Latencies.Mean)
	util.Logf("  99th Percentile Latency: %s", metrics.Latencies.P99)
	util.Logf("  Max Latency: %s", metrics.Latencies.Max)
}

func printTestStatus(metrics *vegeta.Metrics) {
	resourceStatus := util.NewResourceStatus()
	establishedPortsPercent, timeWaitPortsPercent, fdsInUsePercent :=
		resourceStatus.GetPercentages()

	errorMessages := strings.Join(metrics.Errors, ", ")
	if errorMessages == "" {
		errorMessages = "(none)"
	}

	util.Logf(
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
