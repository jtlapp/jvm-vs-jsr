package orderitems

import (
	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type Scenario struct{}

func NewScenario() *Scenario {
	return &Scenario{}
}

func (s *Scenario) GetName() string {
	return "order-items"
}

func (s *Scenario) CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error) {
	backendSetup := database.NewBackendSetup(dbPool, &SetupImpl{dbPool})
	return backendSetup, nil
}

func (s *Scenario) GetTargetProvider(
	baseUrl string,
	randomSeed int64,
	config config.ScenarioConfig,
) func(*vegeta.Target) error {

	trial := NewBenchmarkTrial(baseUrl, randomSeed)
	return func(target *vegeta.Target) error {
		*target = *trial.getRequest()
		return nil
	}
}
