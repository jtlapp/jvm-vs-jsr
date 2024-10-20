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
	PrintMetrics(bs.Metrics)
}

func PrintMetrics(metrics vegeta.Metrics) {
	fmt.Printf("Throughput: %f requests/sec\n", metrics.Throughput)
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Average Latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("99th Percentile Latency: %s\n", metrics.Latencies.P99)
	fmt.Printf("Max Latency: %s\n", metrics.Latencies.Max)
	fmt.Printf("Status Codes: %v\n", metrics.StatusCodes)
}
