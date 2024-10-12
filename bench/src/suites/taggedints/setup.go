package taggedints

import (
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5"
	"jvm-vs-js.jtlapp.com/benchmark/lib"
)

const (
	ROW_COUNT  = 1000000
	MAX_INT    = 1000
	TAG_CHARS  = "0123456789ABCDEF"
	TAG_LENGTH = len(TAG_CHARS)
	SEED       = 12345
)

type SetupImpl struct {
	randGen *rand.Rand
}

func (s *SetupImpl) CreateTables(conn *pgx.Conn) error {
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

func (s *SetupImpl) PopulateTables(conn *pgx.Conn) error {
	for i := 1; i <= ROW_COUNT; i++ {
		tag1 := s.createTag()
		tag2 := s.createTag()
		intVal := s.randGen.Intn(MAX_INT)

		query := `INSERT INTO tagged_ints (tag1, tag2, int, created_at) VALUES ($1, $2, $3, NOW())`
		_, err := conn.Exec(context.Background(), query, tag1, tag2, intVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SetupImpl) GetSharedQueries(conn *pgx.Conn) []lib.SharedQuery {
	return []lib.SharedQuery{
		{
			Name:    "taggedints_sumInts",
			Query:   `SELECT SUM(int) AS sum FROM tagged_ints WHERE tag1=${tag1} AND tag2=${tag2}`,
			Returns: "rows",
		},
		{
			Name:    "taggedints_getInt",
			Query:   `SELECT int FROM tagged_ints WHERE id=${id}`,
			Returns: "rows",
		},
	}
}

// Use the local random generator for tag creation
func (s *SetupImpl) createTag() string {
	return string(TAG_CHARS[s.randGen.Intn(TAG_LENGTH)]) + string(TAG_CHARS[s.randGen.Intn(TAG_LENGTH)])
}
