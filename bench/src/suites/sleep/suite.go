package sleep

import (
	"bytes"
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	sleepDuration = 1000
)

type Suite struct{}

func (s *Suite) GetName() string {
	return "sleep"
}

func (s *Suite) Init() error {
	// nothing to do
	return nil
}

func (s *Suite) SetUpDatabase() error {
	// nothing to do
	return nil
}

func (s *Suite) SetSharedQueries() error {
	// nothing to do
	return nil
}

func (s *Suite) GetTargeter(baseUrl string) vegeta.Targeter {
	url := fmt.Sprintf("%s/api/sleep/%d", baseUrl, sleepDuration)

	return vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    url,
		Body:   bytes.NewBuffer(nil).Bytes(),
	})
}
