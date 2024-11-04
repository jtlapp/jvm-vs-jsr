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
	TotalRequests        int
	Metrics              vegeta.Metrics
	StatusCodes          map[string]int
	ErrorMessages        []string
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
	info *config.PlatformConfig,
	config *config.TestConfig,
	results *TestResults,
	resources *util.ResourceStatus,
) error {
	pool, err := rdb.GetPool()
	if err != nil {
		return err
	}

	appConfigJSON, err := json.Marshal(info.AppConfig)
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
		info.ClientVersion,                         // $1  - clientVersion
		info.AppName,                               // $2  - appName
		info.AppVersion,                            // $3  - appVersion
		appConfigJSON,                              // $4  - appConfig
		info.CPUsPerNode,                           // $5  - cpusPerNode
		clientAction,                               // $6  - clientAction
		config.ScenarioName,                        // $7  - testScenarioName
		config.InitialRequestsPerSecond,            // $8  - testInitialRequestsPerSecond
		config.MaxConnections,                      // $9  - testMaxConnections
		config.WorkerCount,                         // $10 - testWorkerCount
		config.CPUsToUse,                           // $11 - testCPUsUsed
		config.DurationSeconds,                     // $12 - testDurationSeconds
		config.RequestTimeoutSeconds,               // $13 - testTimeoutSeconds
		config.MinWaitSeconds,                      // $14 - testMinWaitSeconds
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
