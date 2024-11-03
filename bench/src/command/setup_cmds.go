package command

import (
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
)

func SetupBackendDB(argsParser *ArgsParser) error {
	backendDB := database.NewBackendDatabase()
	defer backendDB.Close()

	backendSetup, err := createBackendSetup(argsParser, backendDB)
	if err != nil {
		return err
	}
	if err = populateDatabase(backendSetup); err != nil {
		return err
	}
	return assignSharedQueries(backendSetup)
}

func AssignQueries(argsParser *ArgsParser) error {
	backendDB := database.NewBackendDatabase()
	defer backendDB.Close()

	backendSetup, err := createBackendSetup(argsParser, backendDB)
	if err != nil {
		return err
	}
	return assignSharedQueries(backendSetup)
}

func createBackendSetup(argsParser *ArgsParser, backendDB *database.BackendDB) (*database.BackendSetup, error) {
	scenarioName, err := argsParser.GetScenarioName()
	if err != nil {
		return nil, err
	}

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
