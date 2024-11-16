package stats

import (
	"encoding/json"
	"fmt"
	"strings"

	vegeta "github.com/tsenart/vegeta/lib"
)

type TrialInfo struct {
	RandomSeed                   int64
	RequestsPerSecond            float64
	PercentSuccesfullyCompleting float64
	SuccessfulCompletesPerSecond float64
	TotalRequests                uint64
	MeanLatency                  string
	MaxLatency                   string
	Latency50thPercentile        string
	Latency95thPercentile        string
	Latency99thPercentile        string
	Histogram                    interface{}
	StatusCodes                  interface{}
	ErrorMessages                string
}

func NewTrialInfo(metrics *vegeta.Metrics, randomSeed int64) (*TrialInfo, error) {

	histogramJSON, err := json.Marshal(metrics.Histogram)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal histogram: %w", err)
	}
	statusCodesJSON, err := json.Marshal(metrics.StatusCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal status codes: %w", err)
	}

	return &TrialInfo{
		RandomSeed:                   randomSeed,
		RequestsPerSecond:            metrics.Rate,
		PercentSuccesfullyCompleting: metrics.Success,
		SuccessfulCompletesPerSecond: metrics.Throughput,
		TotalRequests:                metrics.Requests,
		MeanLatency:                  metrics.Latencies.Mean.String(),
		MaxLatency:                   metrics.Latencies.Max.String(),
		Latency50thPercentile:        metrics.Latencies.P50.String(),
		Latency95thPercentile:        metrics.Latencies.P95.String(),
		Latency99thPercentile:        metrics.Latencies.P99.String(),
		Histogram:                    histogramJSON,
		StatusCodes:                  statusCodesJSON,
		ErrorMessages:                strings.Join(metrics.Errors, "; "),
	}, nil
}
