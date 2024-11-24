package command

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

var SetupResultsDB = newCommand(
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

var ShowStatistics = newCommand(
	"stats",
	"-scenario=<scenario> [-since=period[d|h|m]] [<trial-options>]",
	"Prints statistics for runs of the given scenario using the given trial "+
		"options. If -since is provided, prints statistics only for trials "+
		"completed since the given time duration.",
	addStatisticsOptions,
	func(commandConfig config.CommandConfig) error {

		platformConfig, err := config.GetPlatformConfig()
		if err != nil {
			return err
		}

		var sinceArg = *commandConfig.SincePeriod
		var sinceDuration time.Duration
		if strings.HasSuffix(sinceArg, "d") {
			days, err := strconv.Atoi(strings.TrimSuffix(sinceArg, "d"))
			if err != nil {
				return fmt.Errorf("failed to parse 'since' period: %v", err)
			}
			sinceDuration = time.Duration(days) * 24 * time.Hour
		} else {
			var err error
			sinceDuration, err = time.ParseDuration(sinceArg)
			if err != nil {
				return fmt.Errorf("failed to parse 'since' period: %v", err)
			}
		}

		startTime := time.Now().Add(-sinceDuration)

		resultsDB := database.NewResultsDatabase()
		defer resultsDB.Close()

		runStats, err := stats.NewRunStats(resultsDB, startTime, platformConfig, &commandConfig)
		if err != nil {
			return err
		}
		util.Log()
		runStats.Print()
		return nil
	})

func addStatisticsOptions(commandConfig *config.CommandConfig, flagSet *flag.FlagSet) {
	commandConfig.SincePeriod = flagSet.String("sincePeriod", "356d",
		"How far back to look for trials. Suffix with 'd' for days, 'h' for hours, or 'm' for minutes.")
	addTrialOptions(commandConfig, flagSet)
}

func confirmWithUser(message string) bool {
	fmt.Printf("%s (Y/n): ", message)
	var response string
	_, _ = fmt.Scanln(&response)
	return response == "Y"
}
