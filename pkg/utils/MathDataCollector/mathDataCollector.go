package MathDataCollector

import (
	"errors"
	"math"
)

type MathDataCollector struct {
	Minimum float64
	Maximum float64
	Sum     float64
	Count   uint64
}

func MakeMathDataCollector() MathDataCollector {
	return MathDataCollector{
		Minimum: math.NaN(),
		Maximum: math.NaN(),
		Sum:     0,
		Count:   0,
	}
}

func NewMathDataCollector() *MathDataCollector {
	result := MakeMathDataCollector()
	return &result
}

func (mathDataCollector *MathDataCollector) Add(data float64) {
	mathDataCollector.Sum += data
	mathDataCollector.Count++
	if math.IsNaN(mathDataCollector.Minimum) {
		mathDataCollector.Minimum = data
	} else {
		mathDataCollector.Minimum = math.Min(data, mathDataCollector.Minimum)
	}
	if math.IsNaN(mathDataCollector.Maximum) {
		mathDataCollector.Maximum = data
	} else {
		mathDataCollector.Maximum = math.Max(data, mathDataCollector.Maximum)
	}
}

func (mathDataCollector *MathDataCollector) GetMinimum() float64 {
	return mathDataCollector.Minimum
}

func (mathDataCollector *MathDataCollector) GetMaximum() float64 {
	return mathDataCollector.Maximum
}

func (mathDataCollector *MathDataCollector) GetSum() float64 {
	return mathDataCollector.Sum
}

func (mathDataCollector *MathDataCollector) GetCount() uint64 {
	return mathDataCollector.Count
}

func (mathDataCollector *MathDataCollector) Average() (float64, error) {
	switch mathDataCollector.Count {
	case 0:
		return math.NaN(), errors.New("can't get average when there is no data")
	default:
		return mathDataCollector.Sum / float64(mathDataCollector.Count), nil
	}
}

