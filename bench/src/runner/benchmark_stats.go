package runner

import (
	vegeta "github.com/tsenart/vegeta/lib"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type BenchmarkStats struct {
	SteadyStateRate int
	Metrics         vegeta.Metrics
}

func (bs BenchmarkStats) Print() {
	util.Log("Steady State Rate: %d", bs.SteadyStateRate)
	PrintMetrics(bs.Metrics)
}

func PrintMetrics(metrics vegeta.Metrics) {
	util.Log("Throughput: %f requests/sec", metrics.Throughput)
	util.Log("Requests: %d", metrics.Requests)
	util.Log("Success Rate: %.2f%%", metrics.Success*100)
	util.Log("Average Latency: %s", metrics.Latencies.Mean)
	util.Log("99th Percentile Latency: %s", metrics.Latencies.P99)
	util.Log("Max Latency: %s", metrics.Latencies.Max)
	util.Log("Status Codes: %v", metrics.StatusCodes)
}
