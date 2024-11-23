package command

import (
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
)

var SetupBackendDB = newCommand(
	"setup-backend",
	"-scenario=<scenario>",
	"Creates database tables and queries required for the test scenario.",
	nil,
	func(clientConfig config.ClientConfig, commandConfig usage.CommandConfig) error {
		backendDB := database.NewBackendDatabase()
		defer backendDB.Close()

		scenarioName := *commandConfig.ScenarioName

		backendSetup, err := createBackendSetup(scenarioName, backendDB)
		if err != nil {
			return err
		}
		if err = populateDatabase(backendSetup); err != nil {
			return err
		}
		return assignSharedQueries(backendSetup)
	})

var AssignQueries = newCommand(
	"assign-queries",
	"-scenario=<scenario>",
	"Sets only the queries required for the test scenario.",
	nil,
	func(clientConfig config.ClientConfig, commandConfig usage.CommandConfig) error {
		backendDB := database.NewBackendDatabase()
		defer backendDB.Close()

		scenarioName := *commandConfig.ScenarioName

		backendSetup, err := createBackendSetup(scenarioName, backendDB)
		if err != nil {
			return err
		}
		return assignSharedQueries(backendSetup)
	})

func createBackendSetup(
	scenarioName string,
	backendDB *database.BackendDB,
) (*database.BackendSetup, error) {

	scenario, err := scenarios.GetScenario(scenarioName)
	if err != nil {
		return nil, err
	}

	dbPool, err := backendDB.GetPool()
	if err != nil {
		return nil, err
	}

	backendSetup, err := scenario.CreateBackendSetup(dbPool)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend setup: %v", err)
	}
	return backendSetup, nil
}

func populateDatabase(backendSetup *database.BackendSetup) error {
	if err := backendSetup.PopulateDatabase(); err != nil {
		return fmt.Errorf("failed to populate database: %v", err)
	}
	return nil
}

func assignSharedQueries(backendSetup *database.BackendSetup) error {
	if err := backendSetup.AssignSharedQueries(); err != nil {
		return fmt.Errorf("failed to assign shared queries: %v", err)
	}
	return nil
}
