package taggedints

import (
	"math/rand"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type Scenario struct {
	databaseSetup *util.DatabaseSetup
}

func (s *Scenario) GetName() string {
	return "taggedints"
}

func (s *Scenario) Init(backendDB *util.BackendDB) error {
	dbPool, err := backendDB.GetPool()
	if err != nil {
		return err
	}
	randGen := rand.New(rand.NewSource(randomSeed))
	impl := &SetupImpl{dbPool, randGen}

	s.databaseSetup = util.NewDatabaseSetup(dbPool, impl)
	return nil
}

func (s *Scenario) SetUpTestTables() error {
	return s.databaseSetup.PopulateDatabase()
}

func (s *Scenario) SetSharedQueries() error {
	return s.databaseSetup.CreateSharedQueries()
}

func (s *Scenario) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	test := NewBenchmarkTest(baseUrl)
	return func(target *vegeta.Target) error {
		*target = *test.getRequest()
		return nil
	}
}
