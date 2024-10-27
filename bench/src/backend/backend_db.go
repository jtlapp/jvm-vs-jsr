package backend

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbUrlEnvVar      = "BACKEND_DATABASE_URL"
	dbUsernameEnvVar = "BACKEND_DATABASE_USERNAME"
	dbPasswordEnvVar = "BACKEND_DATABASE_PASSWORD"
)

type BackendDB struct {
	pool *pgxpool.Pool
}

func NewBackendDatabase() *BackendDB {
	return &BackendDB{}
}

func (db *BackendDB) GetPool() (*pgxpool.Pool, error) {
	if db.pool == nil {
		connConfig, err := pgxpool.ParseConfig(os.Getenv(dbUrlEnvVar))
		if err != nil {
			return nil, err
		}

		connConfig.ConnConfig.User = os.Getenv(dbUsernameEnvVar)
		connConfig.ConnConfig.Password = os.Getenv(dbPasswordEnvVar)

		pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
		db.pool = pool
		return pool, err
	}
	return db.pool, nil
}

func (db *BackendDB) ClosePool() {
	if db.pool != nil {
		db.pool.Close()
	}
}
