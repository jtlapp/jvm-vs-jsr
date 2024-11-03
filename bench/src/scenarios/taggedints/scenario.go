package taggedints

import (
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type Scenario struct{}

func (s *Scenario) GetName() string {
	return "taggedints"
}

func (s *Scenario) CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error) {
	randGen := rand.New(rand.NewSource(randomSeed))
	backendSetup := database.NewBackendSetup(dbPool, &SetupImpl{dbPool, randGen})
	return backendSetup, nil
}

func (s *Scenario) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	test := NewBenchmarkTest(baseUrl)
	return func(target *vegeta.Target) error {
		*target = *test.getRequest()
		return nil
	}
}
