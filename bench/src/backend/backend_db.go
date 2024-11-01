package backend

import (
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

var databaseConfig = util.DatabaseConfig{
	UrlEnvVar:      "BACKEND_DATABASE_URL",
	UsernameEnvVar: "BACKEND_DATABASE_USERNAME",
	PasswordEnvVar: "BACKEND_DATABASE_PASSWORD",
}

type BackendDB struct {
	util.Database
}

func NewBackendDatabase() *BackendDB {
	return &BackendDB{*util.NewDatabase(&databaseConfig)}
}
