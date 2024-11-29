package cmd

import (
	"flag"
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/cli"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

var SetupResultsDB = cli.NewCommand(
	"setup-results",
	"",
	"Creates the results database tables on the client pod.",
	nil,
	func(commandConfig config.CommandConfig) error {
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
			if !confirmWithUser("Drop existing tables?") {
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
	})

var ShowStatistics = cli.NewCommand(
	"stats",
	"-scenario=<scenario> [-count=<num-trials>] [<trial-options>]",
	"Prints statistics for runs of the given scenario using the given trial "+
		"options. If -count is provided, prints statistics only for the "+
		"most recent <num-trials> trials.",
	addStatisticsOptions,
	func(commandConfig config.CommandConfig) error {

		resultsDB := database.NewResultsDatabase()
		defer resultsDB.Close()

		appKeys, err := resultsDB.GetAppKeys()
		if err != nil {
			return err
		}

		if len(appKeys) == 0 {
			fmt.Println("No results found.")
			return nil
		}

		resultSetCount := 0
		for _, appKey := range appKeys {
			runStats, err := stats.NewRunStats(resultsDB, &appKey, &commandConfig,
				*commandConfig.TrialCount)
			if err != nil {
				return err
			}
			if runStats != nil {
				util.Log()
				runStats.Print()
				resultSetCount++
			}
		}
		return nil
	})

func addStatisticsOptions(commandConfig *config.CommandConfig, flagSet *flag.FlagSet) {
	commandConfig.TrialCount = flagSet.Int("count", int((^uint(0))>>1),
		"The number of most recent trials to include in the statistics.")
	addTrialOptions(commandConfig, flagSet)
}

func confirmWithUser(message string) bool {
	fmt.Printf("%s (Y/n): ", message)
	var response string
	_, _ = fmt.Scanln(&response)
	return response == "Y"
}
