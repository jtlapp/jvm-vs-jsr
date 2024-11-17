package singlesleep

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

const (
	sleepDuration = 1000
)

type Scenario struct{}

func NewScenario() *Scenario {
	return &Scenario{}
}

func (s *Scenario) GetName() string {
	return "single-sleep"
}

func (s *Scenario) CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error) {
	return nil, errors.New("this scenario has no backend setup")
}

func (s *Scenario) GetTargetProvider(baseUrl string, randomSeed int64) func(*vegeta.Target) error {
	return func(target *vegeta.Target) error {
		*target = vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/api/sleep/%d", baseUrl, sleepDuration),
			Body:   bytes.NewBuffer(nil).Bytes(),
		}
		return nil
	}
}