package stats

import (
	"time"
)

type LatencyStats struct {
	// Central statistics
	TypicalMedian time.Duration // Average of p50s
	TypicalMean   time.Duration // Average of means

	// Worst case analysis
	WorstMedian time.Duration // Max p50
	WorstP95    time.Duration // Max p95
	WorstP99    time.Duration // Max p99
	WorstMean   time.Duration // Max mean

	// Tail latency analysis
	AverageTailRatio float64 // Average of (p99/p50) ratios
	MaxTailRatio     float64 // Maximum (p99/p50) ratio
}

func CalculateLatencyStats(trials []TrialInfo) (LatencyStats, error) {
	stats := LatencyStats{}

	if len(trials) == 0 {
		return stats, nil
	}

	// Temporary slices to hold converted values
	means := make([]time.Duration, len(trials))
	p50s := make([]time.Duration, len(trials))
	p95s := make([]time.Duration, len(trials))
	p99s := make([]time.Duration, len(trials))
	tailRatios := make([]float64, len(trials))

	// First pass: convert strings to durations and calculate ratios
	var sumMeans, sumP50s time.Duration
	for i, trial := range trials {
		mean, err := time.ParseDuration(trial.MeanLatency)
		if err != nil {
			return stats, err
		}
		means[i] = mean
		sumMeans += mean

		p50, err := time.ParseDuration(trial.Latency50thPercentile)
		if err != nil {
			return stats, err
		}
		p50s[i] = p50
		sumP50s += p50

		p95, err := time.ParseDuration(trial.Latency95thPercentile)
		if err != nil {
			return stats, err
		}
		p95s[i] = p95

		p99, err := time.ParseDuration(trial.Latency99thPercentile)
		if err != nil {
			return stats, err
		}
		p99s[i] = p99

		// Calculate p99/p50 ratio for this trial
		tailRatios[i] = float64(p99) / float64(p50)
	}

	// Calculate averages
	n := float64(len(trials))
	stats.TypicalMean = time.Duration(float64(sumMeans) / n)
	stats.TypicalMedian = time.Duration(float64(sumP50s) / n)

	// Find worst cases
	stats.WorstMean = maxDuration(means)
	stats.WorstMedian = maxDuration(p50s)
	stats.WorstP95 = maxDuration(p95s)
	stats.WorstP99 = maxDuration(p99s)

	// Calculate tail ratio statistics
	stats.MaxTailRatio = maxFloat64(tailRatios)
	stats.AverageTailRatio = averageFloat64(tailRatios)

	return stats, nil
}

func maxDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	max := durations[0]
	for _, d := range durations[1:] {
		if d > max {
			max = d
		}
	}
	return max
}

func maxFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func averageFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var sum float64
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
