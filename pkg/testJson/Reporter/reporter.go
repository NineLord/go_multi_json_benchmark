package Reporter

import "time"

type Reporter struct {
	measurementDuration map[string]int64
}

func NewReporter() *Reporter {
	return &Reporter{
		measurementDuration: make(map[string]int64),
	}
}

func (reporter *Reporter) Measure(functionName string, function func() (any, error)) (any, error) {
	startTime := time.Now().UnixMilli()
	functionResult, functionError := function()
	finishTime := time.Now().UnixMilli()
	duration := finishTime - startTime
	reporter.measurementDuration[functionName] = duration
	return functionResult, functionError
}

func (reporter *Reporter) GetMeasures() map[string]int64 {
	return reporter.measurementDuration
}
