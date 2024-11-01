package results

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"jvm-vs-jsr.jtlapp.com/benchmark/runner"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

var databaseConfig = util.DatabaseConfig{
	UrlEnvVar:      "RESULTS_DATABASE_URL",
	UsernameEnvVar: "RESULTS_DATABASE_USERNAME",
	PasswordEnvVar: "RESULTS_DATABASE_PASSWORD",
}

type TestInfo struct {
	ScenarioName         string
	Action               string
	TestsPerformed       int
	TotalDurationSeconds int
}

type ResultsDB struct {
	util.Database
}

func NewResultsDatabase() *ResultsDB {
	return &ResultsDB{*util.NewDatabase(&databaseConfig)}
}

func (rdb *ResultsDB) SaveResults(
	testInfo *TestInfo,
	config *runner.BenchmarkConfig,
	stats *runner.BenchmarkStats,
	resources *util.ResourceStatus,
) error {
	pool, err := rdb.GetPool()
	if err != nil {
		return err
	}

	appConfigJSON, err := json.Marshal(config.AppConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal app config: %w", err)
	}
	histogramJSON, err := json.Marshal(stats.Metrics.Histogram)
	if err != nil {
		return fmt.Errorf("failed to marshal histogram: %w", err)
	}
	statusCodesJSON, err := json.Marshal(stats.Metrics.StatusCodes)
	if err != nil {
		return fmt.Errorf("failed to marshal status codes: %w", err)
	}

	const query = `
		INSERT INTO results (
			"clientVersion",
			"clientAction",
			"appName",
			"appVersion",
			"appConfig",
			"testCPUsPerNode",
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
			"remainingFDsInUsePercent"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33
		)`

	_, err = pool.Exec(context.Background(), query,
		config.ClientVersion,                     // $1  - clientVersion
		testInfo.Action,                          // $2  - clientAction
		config.AppName,                           // $3  - appName
		config.AppVersion,                        // $4  - appVersion
		appConfigJSON,                            // $5  - appConfig
		config.CPUsPerNode,                       // $6  - testCPUsPerNode
		testInfo.ScenarioName,                    // $7  - testScenarioName
		config.InitialRate,                       // $8  - testInitialRequestsPerSecond
		config.MaxConnections,                    // $9  - testMaxConnections
		config.WorkerCount,                       // $10 - testWorkerCount
		config.CPUsToUse,                         // $11 - testCPUsUsed
		config.DurationSeconds,                   // $12 - testDurationSeconds
		config.RequestTimeoutSeconds,             // $13 - testTimeoutSeconds
		config.MinWaitSeconds,                    // $14 - testMinWaitSeconds
		testInfo.TestsPerformed,                  // $15 - resultTestsPerformed
		testInfo.TotalDurationSeconds,            // $16 - resultTestDurationSeconds
		stats.Metrics.Rate,                       // $17 - resultRequestsPerSecond
		stats.Metrics.Success,                    // $18 - resultSuccessPercent
		stats.Metrics.Throughput,                 // $19 - resultSuccessfulRequestsPerSecond
		stats.Metrics.Requests,                   // $20 - resultTotalRequests
		stats.Metrics.Latencies.Mean.String(),    // $21 - resultMeanLatency
		stats.Metrics.Latencies.P50.String(),     // $22 - result50thPercentileLatency
		stats.Metrics.Latencies.P95.String(),     // $23 - result95thPercentileLatency
		stats.Metrics.Latencies.P99.String(),     // $24 - result99thPercentileLatency
		stats.Metrics.Latencies.Max.String(),     // $25 - resultMaxLatency
		histogramJSON,                            // $26 - resultHistogram
		statusCodesJSON,                          // $27 - resultStatusCodes
		strings.Join(stats.Metrics.Errors, "; "), // $28 - resultErrorMessages
		resources.TotalAvailablePorts,            // $29 - totalAvailablePorts
		resources.TotalFileDescriptors,           // $30 - totalFileDescriptors
		resources.EstablishedPortsCount,          // $31 - remainingPortsActive
		resources.TimeWaitPortsCount,             // $32 - remainingPortsWaiting
		resources.FDsInUseCount,                  // $33 - remainingFDsInUse
	)

	if err != nil {
		return fmt.Errorf("failed to insert results: %w", err)
	}
	return nil
}
