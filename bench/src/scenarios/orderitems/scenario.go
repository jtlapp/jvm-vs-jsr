package orderitems

import (
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type Scenario struct{}

func (s *Scenario) GetName() string {
	return "orderitems"
}

func (s *Scenario) CreateBackendSetup(backendDB *database.BackendDB) (*database.BackendSetup, error) {
	dbPool, err := backendDB.GetPool()
	if err != nil {
		return nil, err
	}
	impl := &SetupImpl{dbPool}

	backendSetup := database.NewBackendSetup(dbPool, impl)
	return backendSetup, nil
}

func (s *Scenario) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	test := NewBenchmarkTest(baseUrl)
	return func(target *vegeta.Target) error {
		*target = *test.getRequest()
		return nil
	}
}
