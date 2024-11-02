package runner

import (
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type Scenario interface {
	GetName() string
	Init(backendDB *database.BackendDB) error
	SetUpTestTables() error
	SetSharedQueries() error
	GetTargetProvider(baseUrl string) func(*vegeta.Target) error
}

type ScenarioFactory func() (Scenario, error)
