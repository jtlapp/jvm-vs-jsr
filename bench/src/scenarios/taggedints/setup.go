package taggedints

import (
	"context"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	totalRows         = 1000000
	maxRandomInt      = 1000
	availableTagChars = "0123456789ABCDEF"
	tagLength         = len(availableTagChars)
	randomSeed        = 12345
)

type SetupImpl struct {
	dbPool  *pgxpool.Pool
	randGen *rand.Rand
}

func (s *SetupImpl) CreateTables() error {
	query := `
        CREATE TABLE IF NOT EXISTS tagged_ints (
          id BIGSERIAL PRIMARY KEY,
          tag1 VARCHAR NOT NULL,
          tag2 VARCHAR NOT NULL,
          int INTEGER NOT NULL,
          created_at TIMESTAMP DEFAULT NOW()
        )`
	_, err := s.dbPool.Exec(context.Background(), query)
	return err
}

func (s *SetupImpl) PopulateTables() error {
	for i := 1; i <= totalRows; i++ {
		tag1 := s.createTag()
		tag2 := s.createTag()
		intVal := s.randGen.Intn(maxRandomInt)

		query := `INSERT INTO tagged_ints (tag1, tag2, int, created_at) VALUES ($1, $2, $3, NOW())`
		_, err := s.dbPool.Exec(context.Background(), query, tag1, tag2, intVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SetupImpl) GetSharedQueries() []util.SharedQuery {
	return []util.SharedQuery{
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
	return string(availableTagChars[s.randGen.Intn(tagLength)]) + string(availableTagChars[s.randGen.Intn(tagLength)])
}
