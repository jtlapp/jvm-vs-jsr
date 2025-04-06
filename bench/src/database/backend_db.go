package database

import (
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

type BackendDB struct {
	Database
}

func NewBackendDatabase() *BackendDB {
	var databaseConfig = DatabaseConfig{
		HostUrlEnvVar:  config.BackendDatabaseUrlEnvVar,
		UsernameEnvVar: config.BackendUsernameEnvVar,
		PasswordEnvVar: config.BackendPasswordEnvVar,
	}
	return &BackendDB{*NewDatabase(&databaseConfig)}
}
