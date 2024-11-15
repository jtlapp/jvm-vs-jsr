package stats

import (
	"math"
	"sort"
)

type ValueStats struct {
	Average float64
	Median  float64
	Lowest  float64
	Highest float64
	Range   float64
	StdDev  float64
	CV      float64 // Coefficient of Variation
}

func CalculateValueStats(values []float64) ValueStats {
	stats := ValueStats{}

	n := len(values)
	if n == 0 {
		return stats
	}

	// Make a copy of the slice to avoid modifying the original
	sorted := make([]float64, n)
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate basic stats
	stats.Lowest = sorted[0]
	stats.Highest = sorted[n-1]
	stats.Range = stats.Highest - stats.Lowest

	// Calculate average
	var sum float64
	for _, v := range values {
		sum += v
	}
	stats.Average = sum / float64(n)

	// Calculate median
	if n%2 == 0 {
		stats.Median = (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		stats.Median = sorted[n/2]
	}

	// Calculate standard deviation and CV
	if n > 1 {
		var sumSquaredDiff float64
		for _, v := range values {
			diff := v - stats.Average
			sumSquaredDiff += diff * diff
		}

		// Using n-1 for sample standard deviation (Bessel's correction)
		variance := sumSquaredDiff / float64(n-1)
		stats.StdDev = math.Sqrt(variance)

		// Calculate CV (only if average is not zero to avoid division by zero)
		if stats.Average != 0 {
			stats.CV = stats.StdDev / stats.Average
		}
	}

	return stats
}
