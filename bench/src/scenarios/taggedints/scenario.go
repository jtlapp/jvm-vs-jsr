package taggedints

import (
	"math/rand"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type Scenario struct {
	backendSetup *database.BackendSetup
}

func (s *Scenario) GetName() string {
	return "taggedints"
}

func (s *Scenario) Init(backendDB *database.BackendDB) error {
	dbPool, err := backendDB.GetPool()
	if err != nil {
		return err
	}
	randGen := rand.New(rand.NewSource(randomSeed))
	impl := &SetupImpl{dbPool, randGen}

	s.backendSetup = database.NewBackendSetup(dbPool, impl)
	return nil
}

func (s *Scenario) SetUpTestTables() error {
	return s.backendSetup.PopulateDatabase()
}

func (s *Scenario) SetSharedQueries() error {
	return s.backendSetup.CreateSharedQueries()
}

func (s *Scenario) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	test := NewBenchmarkTest(baseUrl)
	return func(target *vegeta.Target) error {
		*target = *test.getRequest()
		return nil
	}
}
