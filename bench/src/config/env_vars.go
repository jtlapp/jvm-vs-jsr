package config

import (
	"fmt"
	"os"
)

const (
	BaseAppUrlEnvVar         = "BENCH_BASE_APP_URL"
	BackendDatabaseUrlEnvVar = "BENCH_BACKEND_DB_URL"
	BackendUsernameEnvVar    = "BACKEND_DATABASE_USERNAME"
	BackendPasswordEnvVar    = "BACKEND_DATABASE_PASSWORD"
	MaxReservedPortsEnvVar   = "BENCH_MAX_RESERVED_PORTS"
	ResultsDatabaseUrlEnvVar = "BENCH_RESULTS_DB_URL"
	ResultsUsernameEnvVar    = "RESULTS_DB_USER"
	ResultsPasswordEnvVar    = "RESULTS_DB_PASSWORD"
)

func GetEnvVar(varName string) (string, error) {
	value := os.Getenv(varName)
	if value == "" {
		return "", fmt.Errorf("Missing environment variable: %s", varName)
	}
	return value, nil
}

func GetEnvVarAsUint(varName string) (uint, error) {
	value, err := GetEnvVar(varName)
	if err != nil {
		return 0, err
	}

	var uintValue uint
	_, err = fmt.Sscanf(value, "%d", &uintValue)
	if err != nil {
		return 0, fmt.Errorf("Invalid %s: %s", varName, err)
	}
	return uintValue, nil
}
