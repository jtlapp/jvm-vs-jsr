package taggedints

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

const (
	seed                = 12345
	maxRows             = 1000000
	percentLongRequests = 10
	tagChars            = "0123456789ABCDEF"
	tagCharsLength      = len(tagChars)
)

var loggedResponses = make(map[string]bool)

func main() {
	// Define command-line flags
	rateFlag := flag.Int("rate", 10, "Requests per second")
	durationFlag := flag.Int("duration", 10, "Duration of the test in seconds")

	// Parse the flags
	flag.Parse()

	// Use a local random generator
	r := rand.New(rand.NewSource(seed)) // Replace global rand.Seed

	// Set rate and duration based on command-line input
	rate := vegeta.Rate{Freq: *rateFlag, Per: time.Second}
	duration := time.Duration(*durationFlag) * time.Second

	fmt.Printf("Starting benchmark with rate: %d req/s and duration: %d seconds\n", *rateFlag, *durationFlag)

	targeter := func(tgt *vegeta.Target) error {
		req := getRequest(r) // Pass the local random generator
		*tgt = *req
		return nil
	}

	attacker := vegeta.NewAttacker()

	for res := range attacker.Attack(targeter, rate, duration, "Benchmark Test") {
		body := string(res.Body)
		logResponse(res.Code, body)
	}
}

func getRequest(r *rand.Rand) *vegeta.Target {
	if r.Intn(100) < percentLongRequests {
		return longRequest(r) // Use the local random generator
	}
	return shortRequest(r)
}

func longRequest(r *rand.Rand) *vegeta.Target {
	tag1 := getRandomTag(r)
	tag2 := getRandomTag(r)

	body := fmt.Sprintf(`{"tag1": "%s", "tag2": "%s"}`, tag1, tag2)
	return &vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8080/api/query/taggedints_sumInts",
		Body:   []byte(body),
	}
}

func shortRequest(r *rand.Rand) *vegeta.Target {
	id := r.Intn(maxRows) + 1
	body := fmt.Sprintf(`{"id": %d}`, id)
	return &vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8080/api/query/taggedints_getInt",
		Body:   []byte(body),
	}
}

func getRandomTag(r *rand.Rand) string {
	firstChar := tagChars[r.Intn(tagCharsLength)]
	secondChar := tagChars[r.Intn(tagCharsLength)]
	return string(firstChar) + string(secondChar)
}

func logResponse(status uint16, body string) {
	var comboKey string
	var jsonObj map[string]interface{}

	if body != "" {
		err := json.Unmarshal([]byte(body), &jsonObj)
		query := jsonObj["query"]
		error := jsonObj["error"]

		if err == nil {
			comboKey = fmt.Sprintf("%d|%v|%v", status, query, error)
		} else {
			comboKey = fmt.Sprintf("%d|%s", status, body)
		}
	} else {
		comboKey = fmt.Sprintf("%d", status)
		body = "(empty)"
	}

	if !loggedResponses[comboKey] {
		loggedResponses[comboKey] = true
		fmt.Printf("STATUS %d: %s\n", status, body)
	}
}
