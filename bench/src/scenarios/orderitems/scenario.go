package orderitems

import (
	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type Scenario struct{}

func (s *Scenario) GetName() string {
	return "orderitems"
}

func (s *Scenario) CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error) {
	backendSetup := database.NewBackendSetup(dbPool, &SetupImpl{dbPool})
	return backendSetup, nil
}

func (s *Scenario) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	test := NewBenchmarkTest(baseUrl)
	return func(target *vegeta.Target) error {
		*target = *test.getRequest()
		return nil
	}
}
