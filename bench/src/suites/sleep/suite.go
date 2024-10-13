package sleep

import (
	"bytes"
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-js.jtlapp.com/benchmark/lib"
)

const (
	sleepDuration = 1000
)

type Suite struct{}

func (s *Suite) GetName() string {
	return "sleep"
}

func (s *Suite) Init(backendDB *lib.BackendDB) error {
	// nothing to do
	return nil
}

func (s *Suite) SetUpTestTables() error {
	// nothing to do
	return nil
}

func (s *Suite) SetSharedQueries() error {
	// nothing to do
	return nil
}

func (s *Suite) GetTargetProvider(baseUrl string) func(*vegeta.Target) error {
	return func(target *vegeta.Target) error {
		*target = vegeta.Target{
			Method: "GET",
			URL:    fmt.Sprintf("%s/api/sleep/%d", baseUrl, sleepDuration),
			Body:   bytes.NewBuffer(nil).Bytes(),
		}
		return nil
	}
}
