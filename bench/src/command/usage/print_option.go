package usage

import (
	"fmt"
	"strings"
)

func PrintOption(optionName, argument, description, defaultValue string) {
	if argument != "" {
		argument = "<" + strings.ReplaceAll(argument, " ", "-") + ">"
	}
	fmt.Printf("    -%s %s\n", optionName, argument)
	fmt.Printf("        %s\n", description)
	if defaultValue != "" {
		fmt.Printf("        Default: %s\n", defaultValue)
	}
}
