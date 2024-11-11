package command

import (
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

func SetupResultsDB() error {
	resultsDB := database.NewResultsDatabase()
	defer resultsDB.Close()

	dbPool, err := resultsDB.GetPool()
	if err != nil {
		return err
	}

	tableNames, err := database.GetTableNames(dbPool)
	if err != nil {
		return err
	}
	if len(tableNames) > 0 {
		if !ConfirmWithUser("Drop existing tables?") {
			fmt.Println("Aborted. Database not re-created.")
			fmt.Println()
			return nil
		}
	}

	filter := func(name string) bool { return true }
	err = database.DropTables(dbPool, filter)
	if err != nil {
		return fmt.Errorf("failed to drop tables: %v", err)
	}

	err = resultsDB.CreateTables()
	if err != nil {
		return err
	}

	if len(tableNames) > 0 {
		fmt.Println("Results database re-created.")
	} else {
		fmt.Println("Results database created.")
	}
	fmt.Println()
	return nil
}

func ConfirmWithUser(message string) bool {
	fmt.Printf("%s (Y/n): ", message)
	var response string
	_, _ = fmt.Scanln(&response)
	return response == "Y"
}
