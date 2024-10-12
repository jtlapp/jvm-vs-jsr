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

type DatabaseSetupImpl interface {
	CreateTables(conn *pgx.Conn) error
	PopulateDatabase(conn *pgx.Conn) error
	CreateSharedQueries(conn *pgx.Conn) error
}

type DatabaseSetup struct {
	setupName string
	conn      *pgx.Conn
	impl      DatabaseSetupImpl
}

func CreateDatabaseSetup(setupName string, impl DatabaseSetupImpl) (*DatabaseSetup, error) {
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
	return &DatabaseSetup{setupName, conn, impl}, nil
}

func (bs *DatabaseSetup) GetName() string {
	return bs.setupName
}

func (bs *DatabaseSetup) Run() error {
	if err := DropTables(bs.conn); err != nil {
		return err
	}
	if err := bs.impl.CreateTables(bs.conn); err != nil {
		return err
	}
	if err := bs.impl.PopulateDatabase(bs.conn); err != nil {
		return err
	}
	return bs.impl.CreateSharedQueries(bs.conn)
}

func (bs *DatabaseSetup) RecreateSharedQueries() error {
	if err := EmptyTable(bs.conn, "shared_queries"); err != nil {
		return err
	}
	return bs.impl.CreateSharedQueries(bs.conn)
}

func (bs *DatabaseSetup) Release() error {
	return bs.conn.Close(context.Background())
}
