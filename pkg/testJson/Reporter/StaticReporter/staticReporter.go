package StaticReporter

import (
	"fmt"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter/MeasurementType"
)

var reporter = Reporter.GetInstance()

func Measure[T any](testCount string, jsonName string, measurementType MeasurementType.MeasurementType, function func() (T, error)) (T, error) {
	reporter.StartMeasure(testCount, jsonName, measurementType)
	functionResult, functionError := function()
	if err := reporter.FinishMeasure(testCount, jsonName, measurementType); err != nil {
		panic(fmt.Sprintf("Finished measureting before started: testCount=%s ; jsonName=%s ; measurementType=%v", testCount, jsonName, measurementType))
	}
	return functionResult, functionError
}

func Measure2(testCount string, jsonName string, measurementType MeasurementType.MeasurementType, function func() error) error {
	reporter.StartMeasure(testCount, jsonName, measurementType)
	functionError := function()
	if err := reporter.FinishMeasure(testCount, jsonName, measurementType); err != nil {
		panic(fmt.Sprintf("Finished measureting before started: testCount=%s ; jsonName=%s ; measurementType=%v", testCount, jsonName, measurementType))
	}
	return functionError
}
