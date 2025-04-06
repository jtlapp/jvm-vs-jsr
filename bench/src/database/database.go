package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

type DatabaseConfig struct {
	HostUrlEnvVar  string
	UsernameEnvVar string
	PasswordEnvVar string
}

type Database struct {
	config *DatabaseConfig
	pool   *pgxpool.Pool
}

func NewDatabase(config *DatabaseConfig) *Database {
	return &Database{config: config}
}

func (db *Database) GetPool() (*pgxpool.Pool, error) {
	hostUrl, err := config.GetEnvVar(db.config.HostUrlEnvVar)
	if err != nil {
		return nil, err
	}
	username, err := config.GetEnvVar(db.config.UsernameEnvVar)
	if err != nil {
		return nil, err
	}
	password, err := config.GetEnvVar(db.config.PasswordEnvVar)
	if err != nil {
		return nil, err
	}

	if db.pool == nil {
		connConfig, err := pgxpool.ParseConfig(hostUrl)
		if err != nil {
			return nil, err
		}

		connConfig.ConnConfig.User = username
		connConfig.ConnConfig.Password = password

		pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
		db.pool = pool
		return pool, err
	}
	return db.pool, nil
}

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
