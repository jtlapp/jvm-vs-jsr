package cmd

import (
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/cli"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
)

var ShowAppInfo = cli.NewCommand(
	"app-info",
	"",
	"Prints information about the app, including its load-related configuration.",
	nil,
	func(commandConfig config.CommandConfig) error {

		platformConfig, err := config.GetPlatformConfig()
		if err != nil {
			return err
		}
		platformConfig.Print()
		return nil
	})

var SetupBackendDB = cli.NewCommand(
	"setup-backend",
	"-scenario=<scenario>",
	"Creates database tables and queries required for the test scenario.",
	nil,
	func(commandConfig config.CommandConfig) error {
		backendDB := database.NewBackendDatabase()
		defer backendDB.Close()

		scenarioName := *commandConfig.ScenarioName

		backendSetup, err := createBackendSetup(scenarioName, backendDB)
		if err != nil {
			return err
		}
		return populateDatabase(backendSetup)
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
