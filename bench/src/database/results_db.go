package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const schemaSQL = `
    CREATE TABLE runs (
		"id" SERIAL PRIMARY KEY,
		"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
		"appName" VARCHAR NOT NULL,
		"appVersion" VARCHAR NOT NULL,
		"appConfig" JSONB NOT NULL,
		"cpusPerNode" INTEGER NOT NULL,
		"scenarioName" VARCHAR NOT NULL,
		"initialRequestsPerSecond" INTEGER NOT NULL,
		"initialRandomSeed" INTEGER NOT NULL,
		"maxConnections" INTEGER NOT NULL,
		"workerCount" INTEGER NOT NULL,
		"cpusUsed" INTEGER NOT NULL,
		"trialDurationSeconds" INTEGER NOT NULL,
		"timeoutSeconds" INTEGER NOT NULL,
		"minWaitSeconds" INTEGER NOT NULL,
		"totalRunDurationSeconds" INTEGER NOT NULL,
		"bestTrialID" INTEGER
    );

    CREATE INDEX idx_results_app_name ON runs("appName");
    CREATE INDEX idx_results_app_version ON runs("appVersion");
    CREATE INDEX idx_results_app_config ON runs USING GIN ("appConfig");
    CREATE INDEX idx_results_scenario_name ON runs("scenarioName");

    CREATE TABLE trials (
		"id" SERIAL PRIMARY KEY,
		"createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),
		"runID" INTEGER NOT NULL,
		"randomSeed" INTEGER NOT NULL,
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
	commandConfig *config.CommandConfig,
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
			"appName",
			"appVersion",
			"appConfig",
			"cpusPerNode",
			"scenarioName",
			"initialRequestsPerSecond",
			"initialRandomSeed",
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
		platformConfig.AppName,              // $1  - appName
		platformConfig.AppVersion,           // $2  - appVersion
		appConfigJSON,                       // $3  - appConfig
		platformConfig.CPUsPerNode,          // $4  - cpusPerNode
		commandConfig.ScenarioName,             // $5  - scenarioName
		commandConfig.InitialRequestsPerSecond, // $6  - initialRequestsPerSecond
		commandConfig.InitialRandomSeed,        // $7  - initialRandomSeed
		commandConfig.MaxConnections,           // $8  - maxConnections
		commandConfig.WorkerCount,              // $9 - workerCount
		commandConfig.CPUsToUse,                // $10 - cpusUsed
		commandConfig.DurationSeconds,          // $11 - trialDurationSeconds
		commandConfig.RequestTimeoutSeconds,    // $12 - timeoutSeconds
		commandConfig.MinWaitSeconds,           // $13 - minWaitSeconds
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
	trialInfo *TrialInfo,
	resources *util.ResourceStatus,
) (int, error) {
	pool, err := rdb.GetPool()
	if err != nil {
		return 0, err
	}

	const query = `
		INSERT INTO trials (
			"runID",
			"randomSeed",
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
			$11, $12, $13, $14, $15, $16, $17, $18, $19
		)
		RETURNING id`

	var trialID int
	err = pool.QueryRow(context.Background(), query,
		runID,                                  // $1
		trialInfo.RandomSeed,                   // $2
		trialInfo.RequestsPerSecond,            // $3
		trialInfo.PercentSuccesfullyCompleting, // $4
		trialInfo.SuccessfulCompletesPerSecond, // $5
		trialInfo.TotalRequests,                // $6
		trialInfo.MeanLatency,                  // $7
		trialInfo.MaxLatency,                   // $8
		trialInfo.Latency50thPercentile,        // $9
		trialInfo.Latency95thPercentile,        // $10
		trialInfo.Latency99thPercentile,        // $11
		trialInfo.Histogram,                    // $12
		trialInfo.StatusCodes,                  // $13
		trialInfo.ErrorMessages,                // $14
		resources.TotalAvailablePorts,          // $15
		resources.TotalFileDescriptors,         // $16
		resources.EstablishedPortsCount,        // $17
		resources.TimeWaitPortsCount,           // $18
		resources.FDsInUseCount,                // $19
	).Scan(&trialID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert trial: %w", err)
	}
	return trialID, nil
}

func (rdb *ResultsDB) GetTrials(
	sinceTime time.Time,
	platformConfig *config.PlatformConfig,
	commandConfig *config.CommandConfig,
) ([]TrialInfo, error) {
	pool, err := rdb.GetPool()
	if err != nil {
		return nil, err
	}

	var query = `
		SELECT
			t."randomSeed",
			t."requestsPerSecond",
			t."percentSuccesfullyCompleting",
			t."successfulCompletesPerSecond",
			t."totalRequests",
			t."meanLatency",
			t."maxLatency",
			t."latency50thPercentile",
			t."latency95thPercentile",
			t."latency99thPercentile",
			t."histogram",
			t."statusCodes",
			t."errorMessages"
		FROM trials t
		JOIN runs r ON t.id = r."bestTrialID"
		WHERE r."createdAt" >= $1
		  AND r."appName" = $2
		  AND r."appVersion" = $3
		  AND r."appConfig" = $4
		  AND r."scenarioName" = $5
		  AND r."maxConnections" = $6
		  AND r."workerCount" = $7
		  AND r."cpusUsed" = $8
		  AND r."trialDurationSeconds" = $9
		  AND r."timeoutSeconds" = $10
		  AND r."minWaitSeconds" = $11`

	if *commandConfig.InitialRandomSeed != 0 {
		query += ` AND r."initialRandomSeed" = $12`
	} else {
		// The seed is positive, so this uses all seeds.
		query += ` AND r."initialRandomSeed" >= (-1 * $12)`
	}

	rows, err := pool.Query(context.Background(), query,
		sinceTime,                        // $1  - createdAt
		platformConfig.AppName,           // $2  - appName
		platformConfig.AppVersion,        // $3  - appVersion
		platformConfig.AppConfig,         // $4  - appConfig
		*commandConfig.ScenarioName,          // $5  - scenarioName
		*commandConfig.MaxConnections,        // $6  - maxConnections
		*commandConfig.WorkerCount,           // $7  - workerCount
		*commandConfig.CPUsToUse,             // $8  - cpusUsed
		*commandConfig.DurationSeconds,       // $9 - trialDurationSeconds
		*commandConfig.RequestTimeoutSeconds, // $10 - timeoutSeconds
		*commandConfig.MinWaitSeconds,        // $11 - minWaitSeconds
		*commandConfig.InitialRandomSeed,     // $12 - initialRandomSeed
	)
	if err != nil {
		return nil, err
	}

	var trials []TrialInfo
	for rows.Next() {
		var trial TrialInfo
		err := rows.Scan(
			&trial.RandomSeed,
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
