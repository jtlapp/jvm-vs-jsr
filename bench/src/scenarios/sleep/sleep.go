package sleep

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
)

type SleepScenario struct {
	name string
}

func NewAppSleepScenario() *SleepScenario {
	return &SleepScenario{name: "app-sleep"}
}

func NewPostgresSleepScenario() *SleepScenario {
	return &SleepScenario{name: "pg-sleep"}
}

func (s *SleepScenario) GetName() string {
	return s.name
}

func (s *SleepScenario) CreateBackendSetup(dbPool *pgxpool.Pool) (*database.BackendSetup, error) {
	return nil, errors.New("this scenario has no backend setup")
}

func (s *SleepScenario) GetTargetProvider(
	baseUrl string,
	randomSeed int64,
	config config.ScenarioConfig,
) func(*vegeta.Target) error {

	randGen := rand.New(rand.NewSource(int64(randomSeed)))

	return func(target *vegeta.Target) error {
		var sleepDuration int
		if randGen.Intn(100) < config.PercentLongRequests {
			sleepDuration = config.LongSleepMillis
		} else {
			sleepDuration = config.ShortSleepMillis
		}
		*target = vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/api/%s?millis=%d", baseUrl, s.name, sleepDuration),
			Body:   bytes.NewBuffer(nil).Bytes(),
		}
		return nil
	}
}
