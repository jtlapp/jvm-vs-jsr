package backend

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DropTables(dbPool *pgxpool.Pool) error {
	rows, err := dbPool.Query(context.Background(),
		"SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tablename string
		err = rows.Scan(&tablename)
		if err != nil {
			return err
		}

		if tablename != "shared_queries" {
			query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tablename)
			_, err = dbPool.Exec(context.Background(), query)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func EmptyTable(dbPool *pgxpool.Pool, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s", tableName)
	_, err := dbPool.Exec(context.Background(), query)
	return err
}
