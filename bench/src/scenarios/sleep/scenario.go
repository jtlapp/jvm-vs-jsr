package sleep

import (
	"bytes"
	"errors"
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

const (
	sleepDuration = 1000
)

type Scenario struct{}

func (s *Scenario) GetName() string {
	return "sleep"
}

func (s *Scenario) CreateBackendSetup(backendDB *database.BackendDB) (*database.BackendSetup, error) {
	return nil, errors.New("this scenario has no backend setup")
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
