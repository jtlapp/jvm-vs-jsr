package database

import (
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

type BackendDB struct {
	Database
}

func NewBackendDatabase() *BackendDB {
	var databaseConfig = DatabaseConfig{
		HostEnvVar:         config.HostEnvVar,
		PortEnvVar:         config.PgBouncerPortEnvVar,
		DatabaseNameEnvVar: config.BackendDatabaseNameEnvVar,
		UsernameEnvVar:     config.BackendUsernameEnvVar,
		PasswordEnvVar:     config.BackendPasswordEnvVar,
	}
	return &BackendDB{*NewDatabase(&databaseConfig)}
}
