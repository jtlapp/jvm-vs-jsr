package taggedints

import (
	"fmt"
	"math/rand"
	"net/http"

	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	maxRows             = 1000000
	percentLongRequests = 10
	tagChars            = "0123456789ABCDEF"
	tagCharsLength      = len(tagChars)
)

type BenchmarkTest struct {
	baseUrl string
	randGen *rand.Rand
}

func NewBenchmarkTest(baseUrl string) *BenchmarkTest {
	return &BenchmarkTest{
		baseUrl: baseUrl,
		randGen: rand.New(rand.NewSource(randomSeed)),
	}
}

func (t *BenchmarkTest) getRequest() *vegeta.Target {
	if t.randGen.Intn(100) < percentLongRequests {
		return t.getLongRequest()
	}
	return t.getShortRequest()
}

func (t *BenchmarkTest) getLongRequest() *vegeta.Target {
	tag1 := t.getRandomTag()
	tag2 := t.getRandomTag()

	body := fmt.Sprintf(`{"tag1": "%s", "tag2": "%s"}`, tag1, tag2)
	return &vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/query/taggedints_sumInts", t.baseUrl),
		Header: http.Header{},
		Body:   []byte(body),
	}
}

func (t *BenchmarkTest) getShortRequest() *vegeta.Target {
	id := t.randGen.Intn(maxRows) + 1
	body := fmt.Sprintf(`{"id": %d}`, id)
	return &vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/query/taggedints_getInt", t.baseUrl),
		Header: http.Header{},
		Body:   []byte(body),
	}
}

func (t *BenchmarkTest) getRandomTag() string {
	firstChar := tagChars[t.randGen.Intn(tagCharsLength)]
	secondChar := tagChars[t.randGen.Intn(tagCharsLength)]
	return string(firstChar) + string(secondChar)
}
