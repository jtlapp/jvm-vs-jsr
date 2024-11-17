package command

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/stats"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	sinceTime        = "since"
	defaultSinceTime = "365d"
)

var SetupResultsDB = newCommand(
	"setup-results",
	"",
	"Creates the results database tables on the client pod.",
	nil,
	func(cfg config.ClientConfig) error {
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
	"<scenario> [-since=period[d|h|m]] [<trial-options>]",
	"Prints statistics for runs of the given scenario using the given trial "+
		"options. If -since is provided, prints statistics only for trials "+
		"completed since the given time duration.",
	printStatisticsOptions,
	func(clientConfig config.ClientConfig) error {
		flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		sinceDuration := flagSet.String(sinceTime, defaultSinceTime, "")

		platformConfig, err := config.GetPlatformConfig(clientConfig)
		if err != nil {
			return err
		}
		testConfig, err := getTestConfig(flagSet)
		if err != nil {
			return err
		}

		var sinceTime time.Duration
		if strings.HasSuffix(*sinceDuration, "d") {
			days, err := strconv.Atoi(strings.TrimSuffix(*sinceDuration, "d"))
			if err != nil {
				return fmt.Errorf("failed to parse 'since' period: %v", err)
			}
			sinceTime = time.Duration(days) * 24 * time.Hour
		} else {
			var err error
			sinceTime, err = time.ParseDuration(*sinceDuration)
			if err != nil {
				return fmt.Errorf("failed to parse 'since' period: %v", err)
			}
		}

		startTime := time.Now().Add(-sinceTime)

		resultsDB := database.NewResultsDatabase()
		defer resultsDB.Close()

		runStats, err := stats.NewRunStats(resultsDB, startTime, platformConfig, testConfig)
		if err != nil {
			return err
		}
		util.Log()
		runStats.Print()
		return nil
	})

func printStatisticsOptions() {
	usage.PrintOption(
		sinceTime,
		"since period",
		"How far back to look for trials. Suffix with 'd' for days, 'h' for hours, or 'm' for minutes.",
		defaultSinceTime,
	)
	printTrialOptions()
}

func confirmWithUser(message string) bool {
	fmt.Printf("%s (Y/n): ", message)
	var response string
	_, _ = fmt.Scanln(&response)
	return response == "Y"
}
