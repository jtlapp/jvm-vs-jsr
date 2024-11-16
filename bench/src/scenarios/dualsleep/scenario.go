package dualsleep

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

const (
	longSleepDuration   = 1000
	shortSleepDuration  = 100
	randomSeed          = 12345
	percentLongRequests = 10
)

type Scenario struct{}

func NewScenario() *Scenario {
	return &Scenario{}
}

func (s *Scenario) GetName() string {
	return "dual-sleep"
}

func (s *Scenario) CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error) {
	return nil, errors.New("this scenario has no backend setup")
}

func (s *Scenario) GetTargetProvider(baseUrl string, randonSeed int64) func(*vegeta.Target) error {
	randGen := rand.New(rand.NewSource(randomSeed))

	return func(target *vegeta.Target) error {
		var sleepDuration int
		if randGen.Intn(100) < percentLongRequests {
			sleepDuration = longSleepDuration
		} else {
			sleepDuration = shortSleepDuration
		}
		*target = vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/api/sleep/%d", baseUrl, sleepDuration),
			Body:   bytes.NewBuffer(nil).Bytes(),
		}
		return nil
	}
}
