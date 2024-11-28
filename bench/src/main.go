package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/cli"
	"jvm-vs-jsr.jtlapp.com/benchmark/cmd"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

var commands = []cli.Command{
	cmd.ShowAppConfig,
	cmd.SetupResultsDB,
	cmd.SetupBackendDB,
	cmd.LoopDeterminingRates,
	cmd.DetermineRate,
	cmd.TryRate,
	cmd.ShowStatus,
	cmd.ShowStatistics,
}

func main() {
	framework := cli.Framework{
		Commands:      commands,
		PostParseHook: logBeforeRunning,
		ShowUsage:     showUsage,
		ErrorHook:     logError,
	}
	framework.Run()
}

func logBeforeRunning(flagSet *flag.FlagSet, flagsUsed []string) {
	util.LogOnly("\n========================================")
	commandLine := strings.Join(os.Args, " ")
	util.LogfOnly("\n%s %s", time.Now().Format("2006-01-02 15:04:05"), commandLine)

	configLine := ""
	flagSet.VisitAll(func(f *flag.Flag) {
		if slices.Contains(flagsUsed, f.Name) {
			if configLine != "" {
				configLine += ", "
			}
			configLine += fmt.Sprintf("%s=%s", f.Name, f.Value.String())
		}
	})
	util.Logf("\n  [%s]", configLine)
}

func logError(err error) {
	msg := fmt.Sprintf("Error: %v", err)
	util.Log(msg)
}

func showUsage() {
	fmt.Println()
	fmt.Println("Benchmark tool for testing the performance of a web application.")
	fmt.Println()
	fmt.Println("Commands (use --help to see arguments):")
	for _, cmd := range commands {
		cmd.PrintUsage()
	}
	fmt.Println()
	fmt.Println("Available scenarios:")
	for _, scenario := range scenarios.GetScenarios() {
		fmt.Printf("    %s\n", scenario.GetName())
	}
	fmt.Println()
}
