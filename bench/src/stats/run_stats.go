package stats

import (
	"fmt"
	"time"

	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type RunStats struct {
	TrialCount                   int
	RequestsPerSecond            ValueStats
	SuccessfulCompletesPerSecond ValueStats
	SuccessRate                  ValueStats
	Latency                      LatencyStats
}

func NewRunStats(
	resultsDB *database.ResultsDB,
	startTime time.Time,
	platformConfig *config.PlatformConfig,
	testConfig *config.TestConfig,
) (*RunStats, error) {

	trials, err := resultsDB.GetTrials(startTime, platformConfig, testConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting trials: %v", err)
	}

	if len(trials) == 0 {
		return nil, fmt.Errorf("no runs found meeting these criteria")
	}

	runStats, err := CalculateRunStats(trials)
	if err != nil {
		return nil, fmt.Errorf("error calculating run stats: %v", err)
	}
	return &runStats, nil
}

func CalculateRunStats(trials []database.TrialInfo) (RunStats, error) {
	stats := RunStats{TrialCount: len(trials)}

	// Extract slices for value-based statistics
	rps := make([]float64, len(trials))
	successfulRPS := make([]float64, len(trials))
	successRate := make([]float64, len(trials))

	for i, trial := range trials {
		rps[i] = trial.RequestsPerSecond
		successfulRPS[i] = trial.SuccessfulCompletesPerSecond
		successRate[i] = trial.PercentSuccesfullyCompleting
	}

	// Calculate value-based statistics
	stats.RequestsPerSecond = CalculateValueStats(rps)
	stats.SuccessfulCompletesPerSecond = CalculateValueStats(successfulRPS)
	stats.SuccessRate = CalculateValueStats(successRate)

	// Calculate latency statistics
	latencyStats, err := CalculateLatencyStats(trials)
	if err != nil {
		return stats, err
	}
	stats.Latency = latencyStats

	return stats, nil
}

func (rs RunStats) Print() {
	util.Logf(`Statistics over %d runs:

	Requests Per Second:
		Avg: %.1f, Median: %.1f, Range: %.1f (%.1f to %.1f), SD: %.1f, CV: %.2f

	Successful Completes Per Second:
		Avg: %.1f, Median: %.1f, Range: %.1f (%.1f to %.1f), SD: %.1f, CV: %.2f

	Latency:
		Typical Values:
			Median: %v, Mean: %v
		Worst Cases:
			Median: %v, P95: %v, P99: %v, Mean: %v
		Tail Ratios:
			Average: %.2fx, Max: %.2fx`,
		rs.TrialCount,
		rs.RequestsPerSecond.Average,
		rs.RequestsPerSecond.Median,
		rs.RequestsPerSecond.Range,
		rs.RequestsPerSecond.Lowest,
		rs.RequestsPerSecond.Highest,
		rs.RequestsPerSecond.StdDev,
		rs.RequestsPerSecond.CV,

		rs.SuccessfulCompletesPerSecond.Average,
		rs.SuccessfulCompletesPerSecond.Median,
		rs.SuccessfulCompletesPerSecond.Range,
		rs.SuccessfulCompletesPerSecond.Lowest,
		rs.SuccessfulCompletesPerSecond.Highest,
		rs.SuccessfulCompletesPerSecond.StdDev,
		rs.SuccessfulCompletesPerSecond.CV,

		rs.Latency.TypicalMedian,
		rs.Latency.TypicalMean,
		rs.Latency.WorstMedian,
		rs.Latency.WorstP95,
		rs.Latency.WorstP99,
		rs.Latency.WorstMean,
		rs.Latency.AverageTailRatio,
		rs.Latency.MaxTailRatio)
	util.Log()
}
