package commons

import (
	"cmp"
	"math"
	"slices"

	"github.com/cdipaolo/goml/linear"
	"gonum.org/v1/gonum/stat"
)

const (
	weeksPerYear = 52
)

type rankData struct {
	value float64
	rank float64
}

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
	if len(x) != len(y) {
		Fatalf("x and y must have the same number of elements")
	}
	xMean := Mean(x)
	yMean := Mean(y)
	numerator := 0.0
	denominator1 := 0.0
	denominator2 := 0.0
	for i := range x {
		xDelta := x[i] - xMean
		yDelta := y[i] - yMean
		numerator += xDelta * yDelta
		denominator1 += xDelta * xDelta
		denominator2 += yDelta * yDelta
	}
	denominator := math.Sqrt(denominator1) * math.Sqrt(denominator2)
	correlation := numerator / denominator
	return correlation
}

func GetSpearman(x []float64, y []float64) float64 {
	xRanks := getRanks(x)
	yRanks := getRanks(y)
	correlation := GetCorrelation(xRanks, yRanks)
	return correlation
}

func getRanks(values []float64) []float64 {
	ranks := []rankData{}
	for i, sample := range values {
		data := rankData{
			value: sample,
			rank: float64(i + 1),
		}
		ranks = append(ranks, data)
	}
	slices.SortFunc(ranks, func (a, b rankData) int {
		return cmp.Compare(a.value, b.value)
	})
	output := []float64{}
	for _, data := range ranks {
		output = append(output, data.rank)
	}
	return output
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

func GetRateOfChange(a, b float64) float64 {
	value := a / b - 1.0
	if math.IsNaN(value) || math.IsInf(value, 1) || math.IsInf(value, -1) {
		Fatalf("Invalid ratio: a = %.3f, b = %.3f", a, b)
	}
	return value
}

func GetShortReturn(a, b float64) float64 {
	return - GetRateOfChange(a, b)
}

func GetR2Score(features [][]float64, labels []float64, model *linear.LeastSquares) float64 {
	meanObserved := Mean(labels)
	residualSum := 0.0
	totalSum := 0.0
	for i := range features {
		label := labels[i]
		prediction, err := model.Predict(features[i])
		if err != nil {
			Fatalf("Prediction failed: %v", err)
		}
		residualDelta := label - prediction[0]
		residualSum += residualDelta * residualDelta
		totalDelta := label - meanObserved
		totalSum += totalDelta * totalDelta
	}
	r2Score := 1.0 - residualSum / totalSum
	return r2Score
}