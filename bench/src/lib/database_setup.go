package lib

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// TODO: load via env vars from helm chart; remove env from Dockerfile
	dbURL    = "postgres://pgbouncer-service:6432/testdb"
	username = "user"
	password = "password"
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

func CreateDatabaseSetup(impl DatabaseSetupImpl) (*DatabaseSetup, error) {
	connConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	connConfig.ConnConfig.User = username
	connConfig.ConnConfig.Password = password

	conn, err := pgx.ConnectConfig(context.Background(), connConfig.ConnConfig)
	if err != nil {
		return nil, err
	}
	return &DatabaseSetup{conn, impl}, nil
}

func (bs *DatabaseSetup) PopulateDatabase() error {
	if err := DropTables(bs.conn); err != nil {
		return err
	}
	if err := bs.impl.CreateTables(bs.conn); err != nil {
		return err
	}
	return bs.impl.PopulateTables(bs.conn)
}

func (bs *DatabaseSetup) CreateSharedQueries() error {
	if err := EmptyTable(bs.conn, "shared_queries"); err != nil {
		return err
	}

	sql := `INSERT INTO shared_queries (name, query, returns) VALUES ($1, $2, $3)`

	sharedQueries := bs.impl.GetSharedQueries(bs.conn)
	for _, sharedQuery := range sharedQueries {
		_, err := bs.conn.Exec(context.Background(), sql, sharedQuery.Name, sharedQuery.Query, sharedQuery.Returns)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bs *DatabaseSetup) Release() error {
	return bs.conn.Close(context.Background())
}
