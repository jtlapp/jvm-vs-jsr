package database

import (
	"context"

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
	GetSharedQueries() []SharedQuery
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

func (ds *BackendSetup) AssignSharedQueries() error {
	if err := EmptyTable(ds.pool, "shared_queries"); err != nil {
		return err
	}

	sql := `INSERT INTO shared_queries (name, query, returns) VALUES ($1, $2, $3)`

	sharedQueries := ds.impl.GetSharedQueries()
	for _, sharedQuery := range sharedQueries {
		_, err := ds.pool.Exec(context.Background(), sql, sharedQuery.Name, sharedQuery.Query, sharedQuery.Returns)
		if err != nil {
			return err
		}
	}
	return nil
}
