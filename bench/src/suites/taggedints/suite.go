package taggedints

import (
	"math/rand"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-js.jtlapp.com/benchmark/lib"
)

type Suite struct {
	databaseSetup lib.DatabaseSetup
}

func (s *Suite) GetName() string {
	return "taggedints"
}

func (s *Suite) Init() error {
	impl := &SetupImpl{rand.New(rand.NewSource(SEED))}

	databaseSetup, err := lib.CreateDatabaseSetup(impl)
	if err != nil {
		return err
	}
	s.databaseSetup = *databaseSetup
	return nil
}

func (s *Suite) SetUpDatabase() error {
	return s.databaseSetup.PopulateDatabase()
}

func (s *Suite) SetSharedQueries() error {
	return s.databaseSetup.CreateSharedQueries()
}

func (s *Suite) GetTargeter(baseUrl string) vegeta.Targeter {
	return nil
}
