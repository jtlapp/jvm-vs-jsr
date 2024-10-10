package sleep

import (
	"bytes"
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	sleepDuration = 1000
)

func (s *Suite) GetTargeter(baseUrl string) vegeta.Targeter {
	url := fmt.Sprintf("%s/api/sleep/%d", baseUrl, sleepDuration)

	return vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    url,
		Body:   bytes.NewBuffer(nil).Bytes(),
	})
}
