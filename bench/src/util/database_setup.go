package util

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SharedQuery struct {
	Name    string
	Query   string
	Returns string
}

type DatabaseSetupImpl interface {
	CreateTables() error
	PopulateTables() error
	GetSharedQueries() []SharedQuery
}

type DatabaseSetup struct {
	pool *pgxpool.Pool
	impl DatabaseSetupImpl
}

func NewDatabaseSetup(pool *pgxpool.Pool, impl DatabaseSetupImpl) *DatabaseSetup {
	return &DatabaseSetup{pool, impl}
}

func (ds *DatabaseSetup) PopulateDatabase() error {
	if err := DropTables(ds.pool); err != nil {
		return err
	}
	if err := ds.impl.CreateTables(); err != nil {
		return err
	}
	return ds.impl.PopulateTables()
}

func (ds *DatabaseSetup) CreateSharedQueries() error {
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
