package database

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type TestResults struct {
	TestsPerformed       int
	TotalDurationSeconds int
	Metrics              vegeta.Metrics
}

func (tr *TestResults) Print() {
	metrics := tr.Metrics

	util.Log("Tests Performed: %d", tr.TestsPerformed)
	util.Log("Total Duration: %d seconds", tr.TotalDurationSeconds)
	util.Log("Steady state rate: %.1f", metrics.Rate)
	util.Log("Throughput: %f requests/sec", metrics.Throughput)
	util.Log("Requests: %d", metrics.Requests)
	util.Log("Success Percentage: %.2f%%", metrics.Success*100)
	util.Log("Average Latency: %s", metrics.Latencies.Mean)
	util.Log("99th Percentile Latency: %s", metrics.Latencies.P99)
	util.Log("Max Latency: %s", metrics.Latencies.Max)
	util.Log("Status Codes: %v", metrics.StatusCodes)
}

type ResultsDB struct {
	Database
}

func NewResultsDatabase() *ResultsDB {
	var databaseConfig = DatabaseConfig{
		UrlEnvVar:      "RESULTS_DATABASE_URL",
		UsernameEnvVar: "RESULTS_DATABASE_USERNAME",
		PasswordEnvVar: "RESULTS_DATABASE_PASSWORD",
	}
	return &ResultsDB{*NewDatabase(&databaseConfig)}
}

func (rdb *ResultsDB) SaveResults(
	clientAction string,
	platformConfig *config.PlatformConfig,
	testConfig *config.TestConfig,
	results *TestResults,
	resources *util.ResourceStatus,
) error {
	pool, err := rdb.GetPool()
	if err != nil {
		return err
	}

	appConfigJSON, err := json.Marshal(platformConfig.AppConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal app config: %w", err)
	}
	histogramJSON, err := json.Marshal(results.Metrics.Histogram)
	if err != nil {
		return fmt.Errorf("failed to marshal histogram: %w", err)
	}
	statusCodesJSON, err := json.Marshal(results.Metrics.StatusCodes)
	if err != nil {
		return fmt.Errorf("failed to marshal status codes: %w", err)
	}

	const query = `
		INSERT INTO results (
			"clientVersion",
			"appName",
			"appVersion",
			"appConfig",
			"cpusPerNode",
			"testAction",
			"testScenarioName",
			"testInitialRequestsPerSecond",
			"testMaxConnections",
			"testWorkerCount",
			"testCPUsUsed",
			"testDurationSeconds",
			"testTimeoutSeconds",
			"testMinWaitSeconds",
			"resultTestsPerformed",
			"resultTestDurationSeconds",
			"resultRequestsPerSecond",
			"resultSuccessPercent",
			"resultSuccessfulRequestsPerSecond",
			"resultTotalRequests",
			"resultMeanLatency",
			"result50thPercentileLatency",
			"result95thPercentileLatency",
			"result99thPercentileLatency",
			"resultMaxLatency",
			"resultHistogram",
			"resultStatusCodes",
			"resultErrorMessages",
			"totalAvailablePorts",
			"totalFileDescriptors",
			"remainingPortsActive",
			"remainingPortsWaiting",
			"remainingFDsInUse"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33
		)`

	_, err = pool.Exec(context.Background(), query,
		platformConfig.ClientVersion,               // $1  - clientVersion
		platformConfig.AppName,                     // $2  - appName
		platformConfig.AppVersion,                  // $3  - appVersion
		appConfigJSON,                              // $4  - appConfig
		platformConfig.CPUsPerNode,                 // $5  - cpusPerNode
		clientAction,                               // $6  - clientAction
		testConfig.ScenarioName,                    // $7  - testScenarioName
		testConfig.InitialRequestsPerSecond,        // $8  - testInitialRequestsPerSecond
		testConfig.MaxConnections,                  // $9  - testMaxConnections
		testConfig.WorkerCount,                     // $10 - testWorkerCount
		testConfig.CPUsToUse,                       // $11 - testCPUsUsed
		testConfig.DurationSeconds,                 // $12 - testDurationSeconds
		testConfig.RequestTimeoutSeconds,           // $13 - testTimeoutSeconds
		testConfig.MinWaitSeconds,                  // $14 - testMinWaitSeconds
		results.TestsPerformed,                     // $15 - resultTestsPerformed
		results.TotalDurationSeconds,               // $16 - resultTestDurationSeconds
		results.Metrics.Rate,                       // $17 - resultRequestsPerSecond
		results.Metrics.Success,                    // $18 - resultSuccessPercent
		results.Metrics.Throughput,                 // $19 - resultSuccessfulRequestsPerSecond
		results.Metrics.Requests,                   // $20 - resultTotalRequests
		results.Metrics.Latencies.Mean.String(),    // $21 - resultMeanLatency
		results.Metrics.Latencies.P50.String(),     // $22 - result50thPercentileLatency
		results.Metrics.Latencies.P95.String(),     // $23 - result95thPercentileLatency
		results.Metrics.Latencies.P99.String(),     // $24 - result99thPercentileLatency
		results.Metrics.Latencies.Max.String(),     // $25 - resultMaxLatency
		histogramJSON,                              // $26 - resultHistogram
		statusCodesJSON,                            // $27 - resultStatusCodes
		strings.Join(results.Metrics.Errors, "; "), // $28 - resultErrorMessages
		resources.TotalAvailablePorts,              // $29 - totalAvailablePorts
		resources.TotalFileDescriptors,             // $30 - totalFileDescriptors
		resources.EstablishedPortsCount,            // $31 - remainingPortsActive
		resources.TimeWaitPortsCount,               // $32 - remainingPortsWaiting
		resources.FDsInUseCount,                    // $33 - remainingFDsInUse
	)

	if err != nil {
		return fmt.Errorf("failed to insert results: %w", err)
	}
	return nil
}
