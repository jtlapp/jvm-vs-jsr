package runner

import (
	"encoding/json"
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

const (
	emptyBody = "(empty)"
)

type ResponseLogger struct {
	loggedResponses map[string]bool
}

func NewResponseLogger() *ResponseLogger {
	return &ResponseLogger{
		loggedResponses: make(map[string]bool),
	}
}

func (rl *ResponseLogger) Log(responseCode uint16, body string) {
	var comboKey string
	var jsonObj map[string]interface{}

	if body != "" {
		err := json.Unmarshal([]byte(body), &jsonObj)
		query := jsonObj["query"]
		error := jsonObj["error"]

		if err == nil {
			comboKey = fmt.Sprintf("%d|%v|%v", responseCode, query, error)
		} else {
			comboKey = fmt.Sprintf("%d|%s", responseCode, body)
		}
	} else {
		comboKey = fmt.Sprintf("%d", responseCode)
		body = emptyBody
	}

	if !rl.loggedResponses[comboKey] {
		rl.loggedResponses[comboKey] = true
		if responseCode == 0 && body == emptyBody {
			util.Log("  ex. STATUS: timeout")
		} else {
			util.Log("  ex. STATUS %d: %s", responseCode, body)
		}
	}
}
