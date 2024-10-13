package lib

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbUrlEnvVar      = "DATABASE_URL"
	dbUsernameEnvVar = "DATABASE_USERNAME"
	dbPasswordEnvVar = "DATABASE_PASSWORD"
)

type SharedQuery struct {
	Name    string
	Query   string
	Returns string
}

type DatabaseSetupImpl interface {
	CreateTables(conn *pgx.Conn) error
	PopulateTables(conn *pgx.Conn) error
	GetSharedQueries(conn *pgx.Conn) []SharedQuery
}

type DatabaseSetup struct {
	conn *pgx.Conn
	impl DatabaseSetupImpl
}

func NewDatabaseSetup(impl DatabaseSetupImpl) (*DatabaseSetup, error) {
	connConfig, err := pgxpool.ParseConfig(os.Getenv(dbUrlEnvVar))
	if err != nil {
		return nil, err
	}

	connConfig.ConnConfig.User = os.Getenv(dbUsernameEnvVar)
	connConfig.ConnConfig.Password = os.Getenv(dbPasswordEnvVar)

	conn, err := pgx.ConnectConfig(context.Background(), connConfig.ConnConfig)
	if err != nil {
		return nil, err
	}
	return &DatabaseSetup{conn, impl}, nil
}

func (ds *DatabaseSetup) PopulateDatabase() error {
	if err := DropTables(ds.conn); err != nil {
		return err
	}
	if err := ds.impl.CreateTables(ds.conn); err != nil {
		return err
	}
	return ds.impl.PopulateTables(ds.conn)
}

func (ds *DatabaseSetup) CreateSharedQueries() error {
	if err := EmptyTable(ds.conn, "shared_queries"); err != nil {
		return err
	}

	sql := `INSERT INTO shared_queries (name, query, returns) VALUES ($1, $2, $3)`

	sharedQueries := ds.impl.GetSharedQueries(ds.conn)
	for _, sharedQuery := range sharedQueries {
		_, err := ds.conn.Exec(context.Background(), sql, sharedQuery.Name, sharedQuery.Query, sharedQuery.Returns)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ds *DatabaseSetup) Release() error {
	return ds.conn.Close(context.Background())
}
