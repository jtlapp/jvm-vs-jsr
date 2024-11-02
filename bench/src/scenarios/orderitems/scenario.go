package orderitems

import (
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type Scenario struct {
	backendSetup *database.BackendSetup
}

func (s *Scenario) GetName() string {
	return "orderitems"
}

func (s *Scenario) Init(backendDB *database.BackendDB) error {
	dbPool, err := backendDB.GetPool()
	if err != nil {
		return err
	}
	impl := &SetupImpl{dbPool}

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
