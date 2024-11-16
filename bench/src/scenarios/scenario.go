package scenarios

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/dualsleep"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/orderitems"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/singlesleep"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/taggedints"
)

type Scenario interface {
	GetName() string
	CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error)
	GetTargetProvider(baseUrl string, randomSeed int64) func(*vegeta.Target) error
}

var scenariosSlice = []Scenario{
	singlesleep.NewScenario(),
	dualsleep.NewScenario(),
	taggedints.NewScenario(),
	orderitems.NewScenario(),
}

func GetScenario(name string) (Scenario, error) {
	for _, scenario := range scenariosSlice {
		if scenario.GetName() == name {
			return scenario, nil
		}
	}
	return nil, fmt.Errorf("Scenario not found: %s", name)
}

func GetScenarios() []Scenario {
	return scenariosSlice
}
