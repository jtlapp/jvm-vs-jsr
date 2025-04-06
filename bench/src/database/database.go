package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
)

type DatabaseConfig struct {
	HostEnvVar         string
	PortEnvVar         string
	DatabaseNameEnvVar string
	UsernameEnvVar     string
	PasswordEnvVar     string
}

type Database struct {
	config *DatabaseConfig
	pool   *pgxpool.Pool
}

func NewDatabase(config *DatabaseConfig) *Database {
	return &Database{config: config}
}

func (db *Database) GetPool() (*pgxpool.Pool, error) {
	host, err := config.GetEnvVar(db.config.HostEnvVar)
	if err != nil {
		return nil, err
	}
	port, err := config.GetEnvVar(db.config.PortEnvVar)
	if err != nil {
		return nil, err
	}
	databaseName, err := config.GetEnvVar(db.config.DatabaseNameEnvVar)
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

	url := fmt.Sprintf("postgresql://%s:%s/%s", host, port, databaseName)

	if db.pool == nil {
		connConfig, err := pgxpool.ParseConfig(url)
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
