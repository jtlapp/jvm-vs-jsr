package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const schemaSQL = `
    CREATE TABLE runs (
		"id" SERIAL PRIMARY KEY,
		"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
		"clientVersion" VARCHAR NOT NULL,
		"appName" VARCHAR NOT NULL,
		"appVersion" VARCHAR NOT NULL,
		"appConfig" JSONB NOT NULL,
		"cpusPerNode" INTEGER NOT NULL,
		"scenarioName" VARCHAR NOT NULL,
		"initialRequestsPerSecond" INTEGER NOT NULL,
		"maxConnections" INTEGER NOT NULL,
		"workerCount" INTEGER NOT NULL,
		"cpusUsed" INTEGER NOT NULL,
		"trialDurationSeconds" INTEGER NOT NULL,
		"timeoutSeconds" INTEGER NOT NULL,
		"minWaitSeconds" INTEGER NOT NULL,
		"totalRunDurationSeconds" INTEGER NOT NULL,
		"bestTrialID" INTEGER
    );

    CREATE INDEX idx_results_client_version ON runs("clientVersion");
    CREATE INDEX idx_results_app_name ON runs("appName");
    CREATE INDEX idx_results_app_version ON runs("appVersion");
    CREATE INDEX idx_results_app_config ON runs USING GIN ("appConfig");
    CREATE INDEX idx_results_scenario_name ON runs("scenarioName");

    CREATE TABLE trials (
		"id" SERIAL PRIMARY KEY,
		"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
		"runID" INTEGER NOT NULL,
		"requestsPerSecond" DOUBLE PRECISION NOT NULL,
		"percentSuccesfullyCompleting" DOUBLE PRECISION NOT NULL,
		"successfulCompletesPerSecond" DOUBLE PRECISION NOT NULL,
		"totalRequests" INTEGER NOT NULL,
		"meanLatency" VARCHAR NOT NULL,
		"maxLatency" VARCHAR NOT NULL,
		"latency50thPercentile" VARCHAR NOT NULL,
		"latency95thPercentile" VARCHAR NOT NULL,
		"latency99thPercentile" VARCHAR NOT NULL,
		"histogram" JSONB NOT NULL,
		"statusCodes" JSONB NOT NULL,
		"errorMessages" VARCHAR,
		"availablePorts" INTEGER NOT NULL,
		"fileDescriptors" INTEGER NOT NULL,
		"remainingPortsActive" INTEGER NOT NULL,
		"remainingPortsWaiting" INTEGER NOT NULL,
		"remainingFDsInUse" INTEGER NOT NULL
    );

	ALTER TABLE runs 
		ADD CONSTRAINT fk_runs_best_trial 
		FOREIGN KEY ("bestTrialID") 
		REFERENCES trials(id);

	ALTER TABLE trials 
		ADD CONSTRAINT fk_trials_run 
		FOREIGN KEY ("runID") 
		REFERENCES runs(id);`

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

func (rdb *ResultsDB) CreateTables() error {
	pool, err := rdb.GetPool()
	if err != nil {
		return err
	}

	_, err = pool.Exec(context.Background(), schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	return nil
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
	trialInfo *stats.TrialInfo,
	resources *util.ResourceStatus,
) (int, error) {
	pool, err := rdb.GetPool()
	if err != nil {
		return 0, err
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
		runID,                                  // $1  - runID
		trialInfo.RequestsPerSecond,            // $2
		trialInfo.PercentSuccesfullyCompleting, // $3
		trialInfo.SuccessfulCompletesPerSecond, // $4
		trialInfo.TotalRequests,                // $5
		trialInfo.MeanLatency,                  // $6
		trialInfo.MaxLatency,                   // $7
		trialInfo.Latency50thPercentile,        // $8
		trialInfo.Latency95thPercentile,        // $9
		trialInfo.Latency99thPercentile,        // $10
		trialInfo.Histogram,                    // $11
		trialInfo.StatusCodes,                  // $12
		trialInfo.ErrorMessages,                // $13
		resources.TotalAvailablePorts,          // $14 - availablePorts
		resources.TotalFileDescriptors,         // $15 - fileDescriptors
		resources.EstablishedPortsCount,        // $16 - remainingPortsActive
		resources.TimeWaitPortsCount,           // $17 - remainingPortsWaiting
		resources.FDsInUseCount,                // $18 - remainingFDsInUse
	).Scan(&trialID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert trial: %w", err)
	}
	return trialID, nil
}

func (rdb *ResultsDB) GetTrials(
	sinceTime time.Time,
	platformConfig *config.PlatformConfig,
	testConfig *config.TestConfig,
) ([]stats.TrialInfo, error) {
	pool, err := rdb.GetPool()
	if err != nil {
		return nil, err
	}

	const query = `
		SELECT
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
			"errorMessages"
		FROM trials
		WHERE "createdAt" > $1
		  AND "appName" = $2
		  AND "appVersion" = $3
		  AND "appConfig" = $4
		  AND "scenarioName" = $5
		  AND "initialRequestsPerSecond" = $6
		  AND "maxConnections" = $7
		  AND "workerCount" = $8
		  AND "cpusUsed" = $9
		  AND "trialDurationSeconds" = $10
		  AND "timeoutSeconds" = $11
		  AND "minWaitSeconds" = $12`

	rows, err := pool.Query(context.Background(), query,
		sinceTime,                           // $1  - createdAt
		platformConfig.AppName,              // $2  - appName
		platformConfig.AppVersion,           // $3  - appVersion
		platformConfig.AppConfig,            // $4  - appConfig
		testConfig.ScenarioName,             // $5  - scenarioName
		testConfig.InitialRequestsPerSecond, // $6  - initialRequestsPerSecond
		testConfig.MaxConnections,           // $7  - maxConnections
		testConfig.WorkerCount,              // $8  - workerCount
		testConfig.CPUsToUse,                // $9  - cpusUsed
		testConfig.DurationSeconds,          // $10 - trialDurationSeconds
		testConfig.RequestTimeoutSeconds,    // $11 - timeoutSeconds
		testConfig.MinWaitSeconds,           // $12 - minWaitSeconds
	)
	if err != nil {
		return nil, err
	}

	var trials []stats.TrialInfo
	for rows.Next() {
		var trial stats.TrialInfo
		err := rows.Scan(
			&trial.RequestsPerSecond,
			&trial.PercentSuccesfullyCompleting,
			&trial.SuccessfulCompletesPerSecond,
			&trial.TotalRequests,
			&trial.MeanLatency,
			&trial.MaxLatency,
			&trial.Latency50thPercentile,
			&trial.Latency95thPercentile,
			&trial.Latency99thPercentile,
			&trial.Histogram,
			&trial.StatusCodes,
			&trial.ErrorMessages,
		)
		if err != nil {
			return nil, err
		}
		trials = append(trials, trial)
	}
	return trials, nil
}
