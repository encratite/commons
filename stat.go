package commons

import (
	"gonum.org/v1/gonum/stat"
)

func Mean(samples []float64) float64 {
	return stat.Mean(samples, nil)
}

func StdDev(samples []float64) float64 {
	return stat.StdDev(samples, nil)
}