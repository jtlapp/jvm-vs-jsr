package taggedints

import (
	"context"
	"math/rand"

	"jvm-vs-js.jtlapp.com/benchmark/lib"
)

const (
	ROW_COUNT  = 1000000
	MAX_INT    = 1000
	TAG_CHARS  = "0123456789ABCDEF"
	TAG_LENGTH = len(TAG_CHARS)
	SEED       = 12345
)

type Setup struct {
	lib.DatabaseSetup
	randGen *rand.Rand
}

func (s *Suite) PerformSetup() error {
	baseSetup, err := lib.CreateDatabaseSetup("tagged-ints")
	if err != nil {
		return err
	}
	randGen := rand.New(rand.NewSource(SEED))

	setup := &Setup{*baseSetup, randGen}
	setup.SetActions(setup)
	return setup.Run()
}

func (s *Setup) CreateTables() error {
	query := `
        CREATE TABLE IF NOT EXISTS tagged_ints (
          id BIGSERIAL PRIMARY KEY,
          tag1 VARCHAR NOT NULL,
          tag2 VARCHAR NOT NULL,
          int INTEGER NOT NULL,
          created_at TIMESTAMP DEFAULT NOW()
        )`
	_, err := s.Conn.Exec(context.Background(), query)
	return err
}

func (s *Setup) PopulateDatabase() error {
	for i := 1; i <= ROW_COUNT; i++ {
		tag1 := s.createTag()
		tag2 := s.createTag()
		intVal := s.randGen.Intn(MAX_INT)

		query := `INSERT INTO tagged_ints (tag1, tag2, int, created_at) VALUES ($1, $2, $3, NOW())`
		_, err := s.Conn.Exec(context.Background(), query, tag1, tag2, intVal)
		if err != nil {
			return err 
		}
	}
	return nil
}

func (s *Setup) CreateSharedQueries() error {
	// Replace this with actual shared query implementation.
	return nil
}

// Use the local random generator for tag creation
func (s *Setup) createTag() string {
	return string(TAG_CHARS[s.randGen.Intn(TAG_LENGTH)]) + string(TAG_CHARS[s.randGen.Intn(TAG_LENGTH)])
}
