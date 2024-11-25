package cli

import (
	"flag"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

const (
	fileFlag = "file"
	noFile   = ""
)

type Command interface {
	Name() string
	ArgsUsage() string
	Description() string
	ParseArgs(PostParseHookType) (*config.CommandConfig, error)
	Execute(config.CommandConfig) error
	PrintUsage()
	PrintUsageWithOptions()
}

func NewCommand(
	name, argsUsage, description string,
	addOptions func(*config.CommandConfig, *flag.FlagSet),
	execute func(config.CommandConfig) error,
) Command {
	return &baseCommand{
		name:        name,
		argsUsage:   argsUsage,
		description: description,
		addOptions:  addOptions,
		execute:     execute,
	}
}

type baseCommand struct {
	name        string
	argsUsage   string
	description string
	addOptions  func(*config.CommandConfig, *flag.FlagSet)
	execute     func(config.CommandConfig) error
}

func (c *baseCommand) Name() string        { return c.name }
func (c *baseCommand) ArgsUsage() string   { return c.argsUsage }
func (c *baseCommand) Description() string { return c.description }

func (c *baseCommand) ParseArgs(postParseHook PostParseHookType) (*config.CommandConfig, error) {
	commandConfig := config.CommandConfig{}

	if (*c).addOptions == nil {
		return &commandConfig, nil
	}
	flagSet := flag.NewFlagSet(c.name, flag.ExitOnError)
	(*c).addOptions(&commandConfig, flagSet)
	flagsUsed, err := parseFlagsWithFileDefaults(commandConfig.ConfigFile, flagSet, os.Args[2:])
	if err != nil {
		return nil, err
	}
	if postParseHook != nil {
		postParseHook(flagSet, flagsUsed)
	}
	return &commandConfig, nil
}

func (c *baseCommand) Execute(commandConfig config.CommandConfig) error {
	return c.execute(commandConfig)
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
		commandConfig := config.CommandConfig{}
		c.addOptions(&commandConfig, flagSet)
		flagSet.Usage()
	}
	fmt.Println()
}

func AllowConfigFile(flagSet *flag.FlagSet) *string {
	return flagSet.String(fileFlag, noFile, "Path to YAML config file providing default values.")
}

func installCustomUsageOutput(flagSet *flag.FlagSet) {
	var indent = "    "
	flagSet.Usage = func() {
		flagSet.VisitAll(func(f *flag.Flag) {
			fmt.Printf("%s-%s", indent, f.Name)
			valueDescriptor, flagUsage := flag.UnquoteUsage(f)
			if len(valueDescriptor) > 0 {
				fmt.Print(" " + valueDescriptor)
			}
			fmt.Printf("\n%s%s%s\n", indent, indent, flagUsage)
		})
	}
}

func parseFlagsWithFileDefaults(configFile *string, flagSet *flag.FlagSet, args []string) ([]string, error) {
	flagsUsed := make([]string, 0)

	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	if *configFile == noFile {
		flagSet.VisitAll(func(f *flag.Flag) {
			flagsUsed = append(flagsUsed, f.Name)
		})
	} else {
		providedAsArg := make(map[string]bool)
		flagSet.Visit(func(f *flag.Flag) {
			providedAsArg[f.Name] = true
		})

		configFileBytes, err := os.ReadFile(*configFile)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}

		var configFileMap map[string]interface{}
		if err := yaml.Unmarshal(configFileBytes, &configFileMap); err != nil {
			return nil, fmt.Errorf("error parsing config file: %w", err)
		}

		var fileFlagErr error
		flagSet.VisitAll(func(f *flag.Flag) {
			if f.Name != fileFlag && !providedAsArg[f.Name] {
				if val, ok := configFileMap[f.Name]; ok {
					if err = f.Value.Set(fmt.Sprint(val)); err != nil {
						fileFlagErr = fmt.Errorf("error setting '%s' from file, using default instead", f.Name)
					}
				}
			}
		})
		if fileFlagErr != nil {
			return nil, fileFlagErr
		}

		flagSet.VisitAll(func(f *flag.Flag) {
			if f.Name != fileFlag {
				_, ok := configFileMap[f.Name]
				if providedAsArg[f.Name] || ok {
					flagsUsed = append(flagsUsed, f.Name)
				}
			}
		})
	}
	return flagsUsed, nil
}
