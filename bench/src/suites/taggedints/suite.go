package taggedints

import (
	"math/rand"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-js.jtlapp.com/benchmark/lib"
)

type Suite struct {
	databaseSetup *lib.DatabaseSetup
}

func (s *Suite) GetName() string {
	return "taggedints"
}

func (s *Suite) Init(backendDB *lib.BackendDB) error {
	dbPool, err := backendDB.GetPool()
	if err != nil {
		return err
	}
	randGen := rand.New(rand.NewSource(randomSeed))
	impl := &SetupImpl{dbPool, randGen}

	s.databaseSetup = lib.NewDatabaseSetup(dbPool, impl)
	return nil
}

func (s *Suite) SetUpTestTables() error {
	return s.databaseSetup.PopulateDatabase()
}

func (s *Suite) SetSharedQueries() error {
	return s.databaseSetup.CreateSharedQueries()
}

func (s *Suite) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	test := NewBenchmarkTest(baseUrl)
	return func(target *vegeta.Target) error {
		*target = *test.getRequest()
		return nil
	}
}
