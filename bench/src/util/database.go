package util

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	UrlEnvVar      string
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
	if db.pool == nil {
		connConfig, err := pgxpool.ParseConfig(os.Getenv(db.config.UrlEnvVar))
		if err != nil {
			return nil, err
		}

		connConfig.ConnConfig.User = os.Getenv(db.config.UsernameEnvVar)
		connConfig.ConnConfig.Password = os.Getenv(db.config.PasswordEnvVar)

		pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
		db.pool = pool
		return pool, err
	}
	return db.pool, nil
}

func (db *Database) ClosePool() {
	if db.pool != nil {
		db.pool.Close()
	}
}
