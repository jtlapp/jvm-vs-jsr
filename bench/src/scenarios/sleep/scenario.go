package sleep

import (
	"bytes"
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/backend"
)

const (
	sleepDuration = 1000
)

type Scenario struct{}

func (s *Scenario) GetName() string {
	return "sleep"
}

func (s *Scenario) Init(backendDB *backend.BackendDB) error {
	// nothing to do
	return nil
}

func (s *Scenario) SetUpTestTables() error {
	// nothing to do
	return nil
}

func (s *Scenario) SetSharedQueries() error {
	// nothing to do
	return nil
}

func (s *Scenario) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	return func(target *vegeta.Target) error {
		*target = vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/api/sleep/%d", baseUrl, sleepDuration),
			Body:   bytes.NewBuffer(nil).Bytes(),
		}
		return nil
	}
}
