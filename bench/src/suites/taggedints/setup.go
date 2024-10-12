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

func (s *Suite) PerformSetup() error {
	impl := &SetupImpl{rand.New(rand.NewSource(SEED))}

	databaseSetup, err := lib.CreateDatabaseSetup("tagged-ints", impl)
	if err != nil {
		return err
	}
	databaseSetup.Run()
	return nil
}

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

func (s *SetupImpl) PopulateDatabase(conn *pgx.Conn) error {
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

func (s *SetupImpl) CreateSharedQueries(conn *pgx.Conn) error {
	// Replace this with actual shared query implementation.
	return nil
}

// Use the local random generator for tag creation
func (s *SetupImpl) createTag() string {
	return string(TAG_CHARS[s.randGen.Intn(TAG_LENGTH)]) + string(TAG_CHARS[s.randGen.Intn(TAG_LENGTH)])
}
