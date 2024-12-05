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
	warmupRequestsPerSecond  = 100
	warmupSeconds            = 5
	firstErrorResponseCode   = 400
	timeoutErrorResponseCode = 0
)

type BenchmarkRunner struct {
	platformConfig config.PlatformConfig
	commandConfig  config.CommandConfig
	scenario       *scenarios.Scenario
	resultsDB      *database.ResultsDB
	logger         *ResponseLogger
}

func NewBenchmarkRunner(
	platformConfig config.PlatformConfig,
	commandConfig config.CommandConfig,
	scenario *scenarios.Scenario,
	resultsDB *database.ResultsDB,
) (*BenchmarkRunner, error) {

	return &BenchmarkRunner{
		platformConfig: platformConfig,
		commandConfig:  commandConfig,
		scenario:       scenario,
		resultsDB:      resultsDB,
		logger:         NewResponseLogger(),
	}, nil
}

func (br *BenchmarkRunner) DetermineRate(runCount int, resetRandomSeed bool) (*stats.RunStats, error) {
	if err := br.performWarmupRun(); err != nil {
		return nil, err
	}

	randomSeed := int64(*br.commandConfig.InitialRandomSeed)

	for i := 0; i < runCount; i++ {
		err := br.performRateDetermination(i+1, randomSeed)
		if err != nil {
			return nil, err
		}
		if !resetRandomSeed {
			randomSeed++
		}
	}

	appKey := database.AppKey{
		AppName:    br.platformConfig.AppName,
		AppVersion: br.platformConfig.AppVersion,
		AppConfig:  br.platformConfig.AppConfig,
	}

	return stats.NewRunStats(br.resultsDB, &appKey, &br.commandConfig, runCount)
}

func (br *BenchmarkRunner) TryRate() (metrics *vegeta.Metrics, err error) {
	var runID int

	if err = br.waitForPortsToClear(); err != nil {
		return nil, err
	}
	runID, err = br.resultsDB.CreateRun(&br.platformConfig, &br.commandConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating single-trial run: %v", err)
	}
	util.Log()

	var trialID int
	trialID, metrics, err = br.performRateTrial(
		runID,
		*br.commandConfig.InitialRequestsPerSecond,
		int64(*br.commandConfig.InitialRandomSeed),
		*br.commandConfig.DurationSeconds,
	)
	if err != nil {
		return nil, fmt.Errorf("error performing rate trial: %v", err)
	}

	err = br.resultsDB.UpdateRun(runID, *br.commandConfig.DurationSeconds, trialID)
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
		int64(*br.commandConfig.InitialRandomSeed),
		warmupSeconds,
	)
	if err != nil {
		return fmt.Errorf("error performing warmup trial: %v", err)
	}
	return nil
}

func (br *BenchmarkRunner) performRateDetermination(iteration int, randomSeed int64) error {
	runID, err := br.resultsDB.CreateRun(&br.platformConfig, &br.commandConfig)
	if err != nil {
		return fmt.Errorf("error creating run: %v", err)
	}

	// Find the highest rate of successfully completing requests that the system
	// can handle without any request errors or timeouts.

	requestRateUpperBound := *br.commandConfig.InitialRequestsPerSecond
	requestRateLowerBound := 0
	var lowerBoundMetrics vegeta.Metrics
	currentRequestRate := -1
	nextRequestRate := requestRateUpperBound
	bestTrialID := 0
	var metrics *vegeta.Metrics
	startTime := time.Now()

	util.Logf("\n--- Rate determination run %d ---", iteration)

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
			runID, currentRequestRate, randomSeed, *br.commandConfig.DurationSeconds)
		if err != nil {
			return fmt.Errorf("error performing rate trial: %v", err)
		}
		printTestStatus(metrics)

		if len(metrics.Errors) > 0 || metrics.Throughput < lowerBoundMetrics.Throughput {
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

	if bestTrialID == 0 {
		return fmt.Errorf("something's not right: no successful trials")
	}
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
		br.commandConfig, br.platformConfig.BaseAppUrl, randomSeed)

	attacker := vegeta.NewAttacker(
		vegeta.Workers(uint64(*br.commandConfig.WorkerCount)),
		vegeta.Connections(*br.commandConfig.MaxConnections),
		vegeta.Timeout(time.Duration(*br.commandConfig.RequestTimeoutSeconds)*time.Second),
		vegeta.KeepAlive(true),
	)
	rateLimiter := vegeta.Rate{Freq: rate, Per: time.Second}
	duration := time.Duration(durationSeconds) * time.Second

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targetProvider, rateLimiter, duration, "") {
		br.logger.Log(res.Code, string(res.Body))
		metrics.Add(res)
		if res.Code == timeoutErrorResponseCode || res.Code >= firstErrorResponseCode {
			attacker.Stop()
		}
	}
	attacker.Stop()
	metrics.Close()

	var trialID int
	if runID != 0 {
		resources := util.NewResourceStatus()
		trialInfo, err := database.NewTrialInfo(&metrics, randomSeed)
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
	minDuration := time.Duration(*br.commandConfig.MinSecondsBetweenTests) * time.Second

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

	statusMessage := "100% successful"
	errorMessages := "(none)"
	if len(metrics.Errors) > 0 {
		statusMessage = "Errored"
		errorMessages = strings.Join(metrics.Errors, ", ")
	}

	util.Logf(
		"  %s (%.1f req/s): issued %.1f req/s, %d%% ports active, %d%% ports waiting, %d%% FDs, errors: %s",
		statusMessage,
		metrics.Throughput,
		metrics.Rate,
		uint(establishedPortsPercent+.5),
		uint(timeWaitPortsPercent+.5),
		uint(fdsInUsePercent+.5),
		errorMessages,
	)
}
