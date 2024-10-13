package orderitems

import (
	"fmt"
	"math/rand"

	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	randomSeed     = 123456
	percentUpdates = 50
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
	if t.randGen.Intn(100) < percentUpdates {
		return t.getUpdateRequest()
	}
	return t.getSelectRequest()
}

func (t *BenchmarkTest) getUpdateRequest() *vegeta.Target {
	userNumber := t.randGen.Intn(totalUsers) + 1
	orderNumber := t.randGen.Intn(ordersPerUser) + 1

	orderID := toOrderID(toUserID(userNumber), orderNumber)
	return &vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/query/orderitems_getOrder", t.baseUrl),
		Body:   []byte(fmt.Sprintf(`{"orderID": "%s"}`, orderID)),
	}
}

func (t *BenchmarkTest) getSelectRequest() *vegeta.Target {
	userNumber := t.randGen.Intn(totalUsers) + 1
	orderNumber := t.randGen.Intn(ordersPerUser) + 1

	orderID := toOrderID(toUserID(userNumber), orderNumber)
	return &vegeta.Target{
		Method: "POST",
		URL:    fmt.Sprintf("%s/api/query/orderitems_boostOrderItems", t.baseUrl),
		Body:   []byte(fmt.Sprintf(`{"orderID": "%s"}`, orderID)),
	}
}
