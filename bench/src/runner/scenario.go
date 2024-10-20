package runner

import (
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type Scenario interface {
	GetName() string
	Init(backendDB *util.BackendDB) error
	SetUpTestTables() error
	SetSharedQueries() error
	GetTargetProvider(baseUrl string) func(*vegeta.Target) error
}

type ScenarioFactory func() (Scenario, error)
