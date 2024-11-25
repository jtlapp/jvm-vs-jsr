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

type Framework struct {
	Commands      []Command
	PostParseHook PostParseHookType
	ShowUsage     func()
	ErrorHook     func(error)
}

func (f *Framework) Run() {
	if len(f.Commands) < 1 {
		f.fail(NewUsageError("no command specified"))
	}

	// Extract the command or show help.

	if len(os.Args) == 1 || os.Args[1] == helpOption {
		f.ShowUsage()
		os.Exit(0)
	}
	commandName := os.Args[1]
	command, err := f.find(commandName)
	if err != nil {
		f.fail(err)
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

	commandConfig, err := command.ParseArgs(f.PostParseHook)
	if err != nil {
		f.fail(err)
	}
	err = command.Execute(*commandConfig)
	if err != nil {
		f.fail(err)
	}
}

func (f *Framework) find(name string) (Command, error) {
	for _, c := range f.Commands {
		if c.Name() == name {
			return c, nil
		}
	}
	return nil, NewUsageError("unknown command: %s", name)
}

func (f *Framework) fail(err error) {
	if err != nil {
		if f.ErrorHook != nil {
			f.ErrorHook(err)
		}
		if IsUsageError(err) {
			f.ShowUsage()
		}
		os.Exit(1)
	}
}
