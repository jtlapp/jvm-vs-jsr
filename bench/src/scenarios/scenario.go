package scenarios

import (
	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/command/usage"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/orderitems"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/sleep"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/taggedints"
)

type Scenario interface {
	GetName() string
	CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error)
	GetTargetProvider(
		baseUrl string,
		randomSeed int64,
		config config.ScenarioConfig,
	) func(*vegeta.Target) error
}

var scenariosSlice = []Scenario{
	sleep.NewAppSleepScenario(),
	sleep.NewPostgresSleepScenario(),
	taggedints.NewScenario(),
	orderitems.NewScenario(),
}

func GetScenario(name string) (Scenario, error) {
	if name == "" {
		return nil, usage.NewUsageError("Scenario name is required")
	}
	for _, scenario := range scenariosSlice {
		if scenario.GetName() == name {
			return scenario, nil
		}
	}
	return nil, usage.NewUsageError("Unknown scenario: %s", name)
}

func GetScenarios() []Scenario {
	return scenariosSlice
}
