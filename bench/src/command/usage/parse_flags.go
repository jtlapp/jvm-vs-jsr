package usage

import (
	"flag"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	fileFlag = "file"
	noFile   = ""
)

func ParseFlagsWithFileDefaults(fs *flag.FlagSet, args []string) error {
	configFile := fs.String(fileFlag, noFile, "path to YAML config file providing default values")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	if *configFile != noFile {
		providedAsArg := make(map[string]bool)
		fs.Visit(func(f *flag.Flag) {
			providedAsArg[f.Name] = true
		})

		configFileBytes, err := os.ReadFile(*configFile)
		if err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}

		var configFileMap map[string]interface{}
		if err := yaml.Unmarshal(configFileBytes, &configFileMap); err != nil {
			return fmt.Errorf("error parsing config file: %w", err)
		}

		fs.VisitAll(func(f *flag.Flag) {
			if f.Name != fileFlag && !providedAsArg[f.Name] {
				if val, ok := configFileMap[f.Name]; ok {
					if err = f.Value.Set(fmt.Sprint(val)); err != nil {
						util.Logf("error setting '%s' from file, using default instead", f.Name)
					}
				}
			}
		})

		configString := ""
		fs.VisitAll(func(f *flag.Flag) {
			if f.Name != fileFlag {
				_, ok := configFileMap[f.Name]
				if providedAsArg[f.Name] || ok {
					if configString != "" {
						configString += ", "
					}
					configString += fmt.Sprintf("%s=%s", f.Name, f.Value.String())
				}
			}
		})
		util.Logf("[%s]\n", configString)
	}
	return nil
}
