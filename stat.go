package commons

import (
	"math"
	"slices"

	"gonum.org/v1/gonum/stat"
)

const (
	weeksPerYear = 52
)

func Mean(samples []float64) float64 {
	return stat.Mean(samples, nil)
}

func Median(samples []float64) float64 {
	n := len(samples)
	if n == 0 {
		return math.NaN()
	}
	sortedSamples := make([]float64, n)
	copy(sortedSamples, samples)
	slices.Sort(sortedSamples)
	index := n / 2
	if n % 2 == 0 {
		mean := (sortedSamples[index - 1] + sortedSamples[index]) / 2.0
		return mean
	} else {
		return sortedSamples[index]
	}
}

func StdDev(samples []float64) float64 {
	return stat.StdDev(samples, nil)
}

func GetCorrelation(x []float64, y []float64) float64 {
	xMean := Mean(x)
	yMean := Mean(y)
	numerator := 0.0
	denominator := 0.0
	for i := range x {
		xDelta := x[i] - xMean
		yDelta := y[i] - yMean
		numerator += xDelta * yDelta
		denominator += xDelta * xDelta
	}
	beta := numerator / denominator
	return beta
}

func GetSharpeRatio(weeklyReturns []float64, riskFreeRate float64) float64 {
	if len(weeklyReturns) < 2 {
		return math.NaN()
	}
	meanReturn := Mean(weeklyReturns)
	stdDev := StdDev(weeklyReturns)
	weeklySharpeRatio := (meanReturn - riskFreeRate / weeksPerYear) / stdDev
	sharpeRatio := math.Sqrt(weeksPerYear) * weeklySharpeRatio
	if math.IsInf(sharpeRatio, 1) || math.IsInf(sharpeRatio, -1) {
		return math.NaN()
	}
	return sharpeRatio
}