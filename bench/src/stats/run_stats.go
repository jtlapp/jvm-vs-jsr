package stats

import (
	"fmt"

	"jvm-vs-jsr.jtlapp.com/benchmark/config"
	"jvm-vs-jsr.jtlapp.com/benchmark/database"
	"jvm-vs-jsr.jtlapp.com/benchmark/util"
)

type RunStats struct {
	ScenarioName                 string
	AppKey                       database.AppKey
	TrialCount                   int
	RequestsPerSecond            ValueStats
	SuccessfulCompletesPerSecond ValueStats
	Latency                      LatencyStats
}

func NewRunStats(
	resultsDB *database.ResultsDB,
	appKey *database.AppKey,
	commandConfig *config.CommandConfig,
	maxTrials int,
) (*RunStats, error) {

	trials, err := resultsDB.GetRecentTrials(appKey, commandConfig, maxTrials)
	if err != nil {
		return nil, fmt.Errorf("error getting trials: %v", err)
	}

	if len(trials) == 0 {
		return nil, nil
	}

	runStats, err := CalculateRunStats(*commandConfig.ScenarioName, appKey, trials)
	if err != nil {
		return nil, fmt.Errorf("error calculating run stats: %v", err)
	}
	return &runStats, nil
}

func CalculateRunStats(
	scenarioName string,
	appKey *database.AppKey,
	trials []database.TrialInfo,
) (RunStats, error) {
	stats := RunStats{
		ScenarioName: scenarioName,
		AppKey:       *appKey,
		TrialCount:   len(trials),
	}

	// Extract slices for value-based statistics
	rps := make([]float64, len(trials))
	successfulRPS := make([]float64, len(trials))

	for i, trial := range trials {
		rps[i] = trial.RequestsPerSecond
		successfulRPS[i] = trial.SuccessfulCompletesPerSecond
	}

	// Calculate value-based statistics
	stats.RequestsPerSecond = CalculateValueStats(rps)
	stats.SuccessfulCompletesPerSecond = CalculateValueStats(successfulRPS)

	// Calculate latency statistics
	latencyStats, err := CalculateLatencyStats(trials)
	if err != nil {
		return stats, err
	}
	stats.Latency = latencyStats

	return stats, nil
}

func (rs RunStats) Print() {
	appConfigString, _ := rs.AppKey.AppConfig.ToJsonString()

	util.Logf(`Statistics for %d run(s) of '%s' on %s %s %s:

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
		rs.ScenarioName,
		rs.AppKey.AppName,
		rs.AppKey.AppVersion,
		appConfigString,
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
	util.Log("\n")
}
