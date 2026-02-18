package commons

import (
	"math"
	"slices"

	"gonum.org/v1/gonum/stat"
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