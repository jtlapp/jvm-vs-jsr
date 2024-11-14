package usage

import "os"

func GetScenarioName() (string, error) {
	if len(os.Args) < 3 {
		return "", NewUsageError("scenario name is required")
	}
	return os.Args[2], nil
}
