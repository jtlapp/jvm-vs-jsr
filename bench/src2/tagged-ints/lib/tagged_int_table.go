package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TaggedIntTable struct{}

func CreateTable(conn *pgx.Conn) error {
	query := `
        CREATE TABLE IF NOT EXISTS tagged_ints (
          id BIGSERIAL PRIMARY KEY,
          tag1 VARCHAR NOT NULL,
          tag2 VARCHAR NOT NULL,
          int INTEGER NOT NULL,
          created_at TIMESTAMP DEFAULT NOW()
        )`
	_, err := conn.Exec(context.Background(), query)
	return err
}

func InsertTaggedInt(conn *pgx.Conn, tag1, tag2 string, intVal int) error {
	query := `INSERT INTO tagged_ints (tag1, tag2, int, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := conn.Exec(context.Background(), query, tag1, tag2, intVal)
	return err
}
