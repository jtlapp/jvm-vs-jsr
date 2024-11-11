package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DropTables(dbPool *pgxpool.Pool, filter func(string) bool) error {
	tableNames, err := GetTableNames(dbPool)
	if err != nil {
		return err
	}

	for _, tablename := range tableNames {
		if filter(tablename) {
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

func GetTableNames(dbPool *pgxpool.Pool) ([]string, error) {
	rows, err := dbPool.Query(context.Background(),
		"SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}
	return tableNames, nil
}
