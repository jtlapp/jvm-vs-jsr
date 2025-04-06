package config

import (
	"fmt"
	"os"
)

const (
	AppPortEnvVar             = "APP_PORT"
	BackendDatabaseNameEnvVar = "BACKEND_DB_NAME"
	BackendUsernameEnvVar     = "BACKEND_DATABASE_USERNAME"
	BackendPasswordEnvVar     = "BACKEND_DATABASE_PASSWORD"
	DockerHostEnvVar          = "CONTAINER_HOST"
	HostEnvVar                = "CONTAINER_HOST"
	MaxReservedPortsEnvVar    = "MAX_RESERVED_PORTS"
	PgBouncerPortEnvVar       = "PGBOUNCER_PORT"
	ResultsPortEnvVar         = "RESULTS_DB_PORT"
	ResultsDatabaseNameEnvVar = "RESULTS_DB_NAME"
	ResultsUsernameEnvVar     = "RESULTS_DB_USER"
	ResultsPasswordEnvVar     = "RESULTS_DB_PASSWORD"
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
