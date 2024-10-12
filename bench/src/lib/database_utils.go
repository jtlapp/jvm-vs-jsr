package lib

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func DropTables(conn *pgx.Conn) error {
	rows, err := conn.Query(context.Background(), "SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
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
			_, err = conn.Exec(context.Background(), query)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func EmptyTable(conn *pgx.Conn, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s", tableName)
	_, err := conn.Exec(context.Background(), query)
	return err
}
