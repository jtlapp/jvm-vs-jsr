package runner

import (
	"fmt"

	vegeta "github.com/tsenart/vegeta/lib"
)

type BenchmarkStats struct {
	SteadyStateRate int
	Metrics         vegeta.Metrics
}

func (bs BenchmarkStats) Print() {
	fmt.Printf("Steady State Rate: %d\n", bs.SteadyStateRate)
	fmt.Printf("Requests: %d\n", bs.Metrics.Requests)
	fmt.Printf("Success Rate: %.2f%%\n", bs.Metrics.Success*100)
	fmt.Printf("Average Latency: %s\n", bs.Metrics.Latencies.Mean)
	fmt.Printf("99th Percentile Latency: %s\n", bs.Metrics.Latencies.P99)
	fmt.Printf("Max Latency: %s\n", bs.Metrics.Latencies.Max)
	fmt.Printf("Status Codes: %v\n", bs.Metrics.StatusCodes)
}
