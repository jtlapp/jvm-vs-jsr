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
	util.FLog("Steady State Rate: %d", bs.SteadyStateRate)
	PrintMetrics(bs.Metrics)
}

func PrintMetrics(metrics vegeta.Metrics) {
	util.FLog("Steady state rate: %.1f", metrics.Rate)
	util.FLog("Throughput: %f requests/sec", metrics.Throughput)
	util.FLog("Requests: %d", metrics.Requests)
	util.FLog("Success Percentage: %.2f%%", metrics.Success*100)
	util.FLog("Average Latency: %s", metrics.Latencies.Mean)
	util.FLog("99th Percentile Latency: %s", metrics.Latencies.P99)
	util.FLog("Max Latency: %s", metrics.Latencies.Max)
	util.FLog("Status Codes: %v", metrics.StatusCodes)
}
