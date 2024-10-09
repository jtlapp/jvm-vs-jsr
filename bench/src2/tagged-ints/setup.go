package taggedints

import (
	"math/rand"
)

const (
	ROW_COUNT  = 1000000
	MAX_INT    = 1000
	TAG_CHARS  = "0123456789ABCDEF"
	TAG_LENGTH = len(TAG_CHARS)
	SEED       = 12345
)

type Setup struct {
	BaseSetup
	randGen *rand.Rand
}

func (s *Suite) CreateSetup(dbURL, username, password string) (*Setup, error) {
	baseSetup, err := CreateBaseSetup("tagged-ints", dbURL, username, password)
	if err != nil {
		return nil, err
	}

	randGen := rand.New(rand.NewSource(SEED))

	setup := &Setup{*baseSetup, randGen}
	setup.actions = setup
	
	return setup, nil
}

func (s *Setup) CreateTables() error {
	return CreateTable(s.conn)
}

func (s *Setup) PopulateDatabase() error {
	for i := 1; i <= ROW_COUNT; i++ {
		tag1 := s.createTag()
		tag2 := s.createTag()
		intVal := s.randGen.Intn(MAX_INT)

		err := InsertTaggedInt(s.conn, tag1, tag2, intVal)
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
