package command

import (
	"flag"
	"fmt"
	"os"

	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

type Command interface {
	Name() string
	ArgsUsage() string
	Description() string
	ParseArgs() (*usage.CommandConfig, error)
	Execute(config.ClientConfig, usage.CommandConfig) error
	PrintUsage()
	PrintUsageWithOptions()
}

type baseCommand struct {
	name        string
	argsUsage   string
	description string
	addOptions  func(*usage.CommandConfig, *flag.FlagSet)
	execute     func(config.ClientConfig, usage.CommandConfig) error
}

func (c *baseCommand) Name() string        { return c.name }
func (c *baseCommand) ArgsUsage() string   { return c.argsUsage }
func (c *baseCommand) Description() string { return c.description }

func (c *baseCommand) ParseArgs() (*usage.CommandConfig, error) {
	commandConfig := usage.CommandConfig{}
	if (*c).addOptions != nil {
		flagSet := flag.NewFlagSet(c.name, flag.ExitOnError)
		(*c).addOptions(&commandConfig, flagSet)
		err := usage.ParseFlagsWithFileDefaults(flagSet, os.Args[1:])
		if err != nil {
			return nil, err
		}
	}
	return &commandConfig, nil
}

func (c *baseCommand) Execute(
	clientConfig config.ClientConfig,
	commandConfig usage.CommandConfig,
) error {
	return c.execute(clientConfig, commandConfig)
}

func (c *baseCommand) PrintUsage() {
	fmt.Printf("    %s %s\n", c.Name(), c.ArgsUsage())
	fmt.Printf("        %s\n", c.Description())
}

func (c *baseCommand) PrintUsageWithOptions() {
	fmt.Println("Usage:")
	c.PrintUsage()

	if c.addOptions != nil {
		fmt.Println("Options:")
		flagSet := flag.NewFlagSet(c.Name(), flag.ExitOnError)
		installCustomUsageOutput(flagSet)
		commandConfig := usage.CommandConfig{}
		c.addOptions(&commandConfig, flagSet)
		flagSet.Usage()
	}
	fmt.Println()
}

func newCommand(
	name, argsUsage, description string,
	addOptions func(*usage.CommandConfig, *flag.FlagSet),
	execute func(config.ClientConfig, usage.CommandConfig) error,
) Command {

	return &baseCommand{
		name:        name,
		argsUsage:   argsUsage,
		description: description,
		addOptions:  addOptions,
		execute:     execute,
	}
}

var Commands = []Command{
	SetupResultsDB,
	SetupBackendDB,
	AssignQueries,
	LoopDeterminingRates,
	DetermineRate,
	TryRate,
	ShowStatus,
	ShowStatistics,
}

func Find(name string) (Command, error) {
	for _, c := range Commands {
		if c.Name() == name {
			return c, nil
		}
	}
	return nil, usage.NewUsageError("unknown command: %s", name)
}

func installCustomUsageOutput(fs *flag.FlagSet) {
	var indent = "    "
    fs.Usage = func() {
        fs.VisitAll(func(f *flag.Flag) {
			fmt.Printf("%s-%s", indent, f.Name)
			valueDescriptor, flagUsage := flag.UnquoteUsage(f)
			if len(valueDescriptor) > 0 {
				fmt.Print(" " + valueDescriptor)
			}
			fmt.Printf("\n%s%s%s\n", indent, indent, flagUsage)
        })
    }
}
