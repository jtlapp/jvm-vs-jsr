package cli

import (
	"flag"
	"os"
	"slices"
	"strings"
)

const (
	helpOption = "-help"
)

type PostParseHookType func(flagSet *flag.FlagSet, flagsUsed []string)

type Runner struct {
	Commands      []Command
	PostParseHook PostParseHookType
	ShowUsage     func()
	ErrorHook     func(error)
}

func (r *Runner) Run() {
	if len(r.Commands) < 1 {
		r.fail(NewUsageError("no command specified"))
	}

	// Extract the command or show help.

	if len(os.Args) == 1 || os.Args[1] == helpOption {
		r.ShowUsage()
		os.Exit(0)
	}
	commandName := os.Args[1]
	command, err := r.find(commandName)
	if err != nil {
		r.fail(err)
	}

	// Show command-specific help if requested.

	index := slices.IndexFunc(os.Args, func(arg string) bool {
		return strings.HasSuffix(arg, helpOption)
	})
	if index != -1 {
		command.PrintUsageWithOptions()
		os.Exit(0)
	}

	// Execute the command.

	commandConfig, err := command.ParseArgs(r.PostParseHook)
	if err != nil {
		r.fail(err)
	}
	err = command.Execute(*commandConfig)
	if err != nil {
		r.fail(err)
	}
}

func (r *Runner) find(name string) (Command, error) {
	for _, c := range r.Commands {
		if c.Name() == name {
			return c, nil
		}
	}
	return nil, NewUsageError("unknown command: %s", name)
}

func (c *Runner) fail(err error) {
	if err != nil {
		if c.ErrorHook != nil {
			c.ErrorHook(err)
		}
		if IsUsageError(err) {
			c.ShowUsage()
		}
		os.Exit(1)
	}
}
