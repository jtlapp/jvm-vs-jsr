package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type SharedQuery struct {
	Name    string
	Query   string
	Returns string
}

type BackendSetupImpl interface {
	CreateTables() error
	PopulateTables() error
}

type BackendSetup struct {
	pool *pgxpool.Pool
	impl BackendSetupImpl
}

func NewBackendSetup(pool *pgxpool.Pool, impl BackendSetupImpl) *BackendSetup {
	return &BackendSetup{pool, impl}
}

func (ds *BackendSetup) PopulateDatabase() error {
	filter := func(name string) bool { return name != "shared_queries" }

	if err := DropTables(ds.pool, filter); err != nil {
		return err
	}
	if err := ds.impl.CreateTables(); err != nil {
		return err
	}
	return ds.impl.PopulateTables()
}
