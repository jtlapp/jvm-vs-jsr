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

type DatabaseSetupActions interface {
	CreateTables() error
	PopulateDatabase() error
	CreateSharedQueries() error
}

type DatabaseSetup struct {
	setupName string
	Conn      *pgx.Conn
	actions   DatabaseSetupActions
}

func CreateDatabaseSetup(setupName string) (*DatabaseSetup, error) {
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

	return &DatabaseSetup{setupName: setupName, Conn: conn}, nil
}

func (bs *DatabaseSetup) SetActions(actions DatabaseSetupActions) {
	bs.actions = actions
}

func (bs *DatabaseSetup) GetName() string {
	return bs.setupName
}

func (bs *DatabaseSetup) Run() error {
	if err := DropTables(bs.Conn); err != nil {
		return err
	}
	if err := bs.actions.CreateTables(); err != nil {
		return err
	}
	if err := bs.actions.PopulateDatabase(); err != nil {
		return err
	}
	return bs.actions.CreateSharedQueries()
}

func (bs *DatabaseSetup) RecreateSharedQueries() error {
	if err := EmptyTable(bs.Conn, "shared_queries"); err != nil {
		return err
	}
	return bs.actions.CreateSharedQueries()
}

func (bs *DatabaseSetup) Release() error {
	return bs.Conn.Close(context.Background())
}
