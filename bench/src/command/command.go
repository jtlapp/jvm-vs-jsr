package command

import (
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

type Command interface {
	Name() string
	ArgsUsage() string
	Description() string
	Execute(clientConfig config.ClientConfig) error
	PrintUsage()
	PrintUsageWithOptions()
}

type baseCommand struct {
	name         string
	argsUsage    string
	description  string
	printOptions func()
	execute      func(config.ClientConfig) error
}

func (c *baseCommand) Name() string                          { return c.name }
func (c *baseCommand) ArgsUsage() string                     { return c.argsUsage }
func (c *baseCommand) Description() string                   { return c.description }
func (c *baseCommand) Execute(cfg config.ClientConfig) error { return c.execute(cfg) }

func (c *baseCommand) PrintUsage() {
	fmt.Printf("    %s %s\n", c.Name(), c.ArgsUsage())
	fmt.Printf("        %s\n", c.Description())
}

func (c *baseCommand) PrintUsageWithOptions() {
	fmt.Println("Usage:")
	c.PrintUsage()
	if c.printOptions != nil {
		fmt.Println("Options:")
		c.printOptions()
	}
	fmt.Println()
}

func newCommand(
	name, argsUsage, description string,
	printOptions func(),
	execute func(config.ClientConfig) error) Command {

	return &baseCommand{
		name:         name,
		argsUsage:    argsUsage,
		description:  description,
		printOptions: printOptions,
		execute:      execute,
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
}

func Find(name string) (Command, error) {
	for _, c := range Commands {
		if c.Name() == name {
			return c, nil
		}
	}
	return nil, usage.NewUsageError("unknown command: %s", name)
}
