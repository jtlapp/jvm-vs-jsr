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

func (rdb *ResultsDB) CreateRun(
	platformConfig *config.PlatformConfig,
	testConfig *config.TestConfig,
) (int, error) {
	pool, err := rdb.GetPool()
	if err != nil {
		return 0, err
	}

	appConfigJSON, err := json.Marshal(platformConfig.AppConfig)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal app config: %w", err)
	}

	const query = `
		INSERT INTO runs (
			"clientVersion",
			"appName",
			"appVersion",
			"appConfig",
			"cpusPerNode",
			"scenarioName",
			"initialRequestsPerSecond",
			"maxConnections",
			"workerCount",
			"cpusUsed",
			"trialDurationSeconds",
			"timeoutSeconds",
			"minWaitSeconds",
			"totalRunDurationSeconds"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		RETURNING id`

	var runID int
	err = pool.QueryRow(context.Background(), query,
		platformConfig.ClientVersion,        // $1  - clientVersion
		platformConfig.AppName,              // $2  - appName
		platformConfig.AppVersion,           // $3  - appVersion
		appConfigJSON,                       // $4  - appConfig
		platformConfig.CPUsPerNode,          // $5  - cpusPerNode
		testConfig.ScenarioName,             // $6  - scenarioName
		testConfig.InitialRequestsPerSecond, // $7  - initialRequestsPerSecond
		testConfig.MaxConnections,           // $8  - maxConnections
		testConfig.WorkerCount,              // $9  - workerCount
		testConfig.CPUsToUse,                // $10 - cpusUsed
		testConfig.DurationSeconds,          // $11 - trialDurationSeconds
		testConfig.RequestTimeoutSeconds,    // $12 - timeoutSeconds
		testConfig.MinWaitSeconds,           // $13 - minWaitSeconds
		0,                                   // $14 - totalRunDurationSeconds (initialized to 0)
	).Scan(&runID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert run: %w", err)
	}
	return runID, nil
}

func (rdb *ResultsDB) UpdateRun(runID int, totalDurationSeconds int, bestTrialID int) error {
	pool, err := rdb.GetPool()
	if err != nil {
		return err
	}

	const query = `
		UPDATE runs 
		SET "totalRunDurationSeconds" = $1,
		    "bestTrialID" = $2
		WHERE id = $3`

	_, err = pool.Exec(context.Background(), query,
		totalDurationSeconds, // $1
		bestTrialID,          // $2
		runID,                // $3
	)

	if err != nil {
		return fmt.Errorf("failed to update run: %w", err)
	}
	return nil
}

func (rdb *ResultsDB) SaveTrial(
	runID int,
	metrics *vegeta.Metrics,
	resources *util.ResourceStatus,
) (int, error) {
	pool, err := rdb.GetPool()
	if err != nil {
		return 0, err
	}

	histogramJSON, err := json.Marshal(metrics.Histogram)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal histogram: %w", err)
	}
	statusCodesJSON, err := json.Marshal(metrics.StatusCodes)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal status codes: %w", err)
	}

	const query = `
		INSERT INTO trials (
			"runID",
			"requestsPerSecond",
			"percentSuccesfullyCompleting",
			"successfulCompletesPerSecond",
			"totalRequests",
			"meanLatency",
			"maxLatency",
			"latency50thPercentile",
			"latency95thPercentile",
			"latency99thPercentile",
			"histogram",
			"statusCodes",
			"errorMessages",
			"availablePorts",
			"fileDescriptors",
			"remainingPortsActive",
			"remainingPortsWaiting",
			"remainingFDsInUse"
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18
		)
		RETURNING id`

	var trialID int
	err = pool.QueryRow(context.Background(), query,
		runID,                             // $1  - runID
		metrics.Rate,                      // $2  - requestsPerSecond
		metrics.Success,                   // $3  - percentSuccesfullyCompleting
		metrics.Throughput,                // $4  - successfulCompletesPerSecond
		metrics.Requests,                  // $5  - totalRequests
		metrics.Latencies.Mean.String(),   // $6  - meanLatency
		metrics.Latencies.Max.String(),    // $7  - maxLatency
		metrics.Latencies.P50.String(),    // $8  - latency50thPercentile
		metrics.Latencies.P95.String(),    // $9  - latency95thPercentile
		metrics.Latencies.P99.String(),    // $10 - latency99thPercentile
		histogramJSON,                     // $11 - histogram
		statusCodesJSON,                   // $12 - statusCodes
		strings.Join(metrics.Errors, ";"), // $13 - errorMessages
		resources.TotalAvailablePorts,     // $14 - availablePorts
		resources.TotalFileDescriptors,    // $15 - fileDescriptors
		resources.EstablishedPortsCount,   // $16 - remainingPortsActive
		resources.TimeWaitPortsCount,      // $17 - remainingPortsWaiting
		resources.FDsInUseCount,           // $18 - remainingFDsInUse
	).Scan(&trialID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert trial: %w", err)
	}
	return trialID, nil
}
