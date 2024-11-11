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

func (br *BenchmarkRunner) DetermineRate() (*vegeta.Metrics, error) {
	if err := br.waitForPortsToClear(); err != nil {
		return nil, err
	}
	runID, err := br.resultsDB.CreateRun(&br.platformConfig, &br.testConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating run: %v", err)
	}

	// Warm up the application, in case it does JIT.

	util.Log()
	util.Log("Warmup run (ignored)...")
	_, _, err = br.performRateTrial(0, warmupRequestsPerSecond, warmupSeconds)
	if err != nil {
		return nil, fmt.Errorf("error performing warmup trial: %v", err)
	}

	// Find the highest rate that the system can handle without errors.

	rateUpperBound := br.testConfig.InitialRequestsPerSecond
	rateLowerBound := 0
	currentRate := -1
	nextRate := rateUpperBound
	testsPerformed := 0
	var bestTrialID int
	var metrics *vegeta.Metrics
	startTime := time.Now()

	for currentRate != 0 && nextRate != currentRate {
		br.waitBetweenTests()

		currentRate = nextRate
		if currentRate == rateLowerBound {
			break
		}

		util.Log()
		util.FLog("Testing %d requests/sec...", currentRate)
		bestTrialID, metrics, err = br.performRateTrial(runID, currentRate, br.testConfig.DurationSeconds)
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

	totalDurationSeconds := int(time.Since(startTime).Seconds())
	err = br.resultsDB.UpdateRun(runID, totalDurationSeconds, bestTrialID)
	if err != nil {
		return nil, fmt.Errorf("error updating multi-trial run: %v", err)
	}

	return metrics, nil
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
	trialID, metrics, err = br.performRateTrial(runID, br.testConfig.InitialRequestsPerSecond, br.testConfig.DurationSeconds)
	if err != nil {
		return nil, fmt.Errorf("error performing rate trial: %v", err)
	}

	err = br.resultsDB.UpdateRun(runID, br.testConfig.DurationSeconds, trialID)
	if err != nil {
		return nil, fmt.Errorf("error updating single-trial run: %v", err)
	}

	return metrics, err
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

func (br *BenchmarkRunner) performRateTrial(runID, rate, durationSeconds int) (int, *vegeta.Metrics, error) {

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
	attacker.Stop()
	metrics.Close()

	var trialID int
	var err error
	if runID != 0 {
		resources := util.NewResourceStatus()
		trialID, err = br.resultsDB.SaveTrial(runID, &metrics, &resources)
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
