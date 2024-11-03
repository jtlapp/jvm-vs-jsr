package scenarios

import (
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/orderitems"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/sleep"
	"jvm-vs-jsr.jtlapp.com/benchmark/scenarios/taggedints"
)

type Scenario interface {
	GetName() string
	CreateBackendSetup(backendDB *database.BackendDB) (*database.BackendSetup, error)
	GetTargetProvider(baseUrl string) func(*vegeta.Target) error
}

var scenariosSlice = []Scenario{
	&sleep.Scenario{},
	&taggedints.Scenario{},
	&orderitems.Scenario{},
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
