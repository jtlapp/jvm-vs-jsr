package taggedints

import (
	"math/rand"

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
	return "tagged-ints"
}

func (s *Scenario) CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error) {
	randGen := rand.New(rand.NewSource(randomSeed))
	backendSetup := database.NewBackendSetup(dbPool, &SetupImpl{dbPool, randGen})
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
