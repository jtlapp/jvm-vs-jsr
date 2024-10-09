package lib

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SetupActions interface {
	CreateTables() error
	PopulateDatabase() error
	CreateSharedQueries() error
}

type BaseSetup struct {
	setupName string
	conn      *pgx.Conn
	actions   SetupActions
}

func CreateBaseSetup(setupName, dbURL, username, password string) (*BaseSetup, error) {
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

	return &BaseSetup{setupName: setupName, conn: conn}, nil
}

func (bs *BaseSetup) GetName() string {
	return bs.setupName
}

func (bs *BaseSetup) Run() error {
	if err := DropTables(bs.conn); err != nil {
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

func (bs *BaseSetup) RecreateSharedQueries() error {
	if err := EmptyTable(bs.conn, "shared_queries"); err != nil {
		return err
	}
	return bs.actions.CreateSharedQueries()
}

func (bs *BaseSetup) Release() error {
	return bs.conn.Close(context.Background())
}
