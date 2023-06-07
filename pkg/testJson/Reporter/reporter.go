package Reporter

import (
	"fmt"
	"sync"
)

var lock = &sync.Mutex{}
var instance *Reporter

func GetInstance() *Reporter {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = newReporter()
		}
	}

	return instance
}

type ReportData map[string]map[string]map[MeasurementType]*Measurement

type Reporter struct {
	lock                sync.RWMutex
	measurementDuration ReportData
}

func newReporter() *Reporter {
	return &Reporter{
		lock:                sync.RWMutex{},
		measurementDuration: make(ReportData),
	}
}

func (reporter *Reporter) StartMeasure(testCount *string, jsonName *string, measurementType MeasurementType) {
	reporter.lock.Lock()
	defer reporter.lock.Unlock()
	testMap := setDefault(reporter.measurementDuration, *testCount, make(map[string]map[MeasurementType]*Measurement))
	jsonMap := setDefault(testMap, *jsonName, make(map[MeasurementType]*Measurement))
	jsonMap[measurementType] = NewMeasurement()
}

func (reporter *Reporter) FinishMeasure(testCount *string, jsonName *string, measurementType MeasurementType) error {
	reporter.lock.Lock()
	defer reporter.lock.Unlock()

	if testMap, ok := reporter.measurementDuration[*testCount]; !ok {
		return fmt.Errorf("can't find test count: %s", *testCount)
	} else if jsonMap, ok := testMap[*jsonName]; !ok {
		return fmt.Errorf("can't find json name: %s", *jsonName)
	} else if measurement, ok := jsonMap[measurementType]; !ok {
		return fmt.Errorf("can't find measurement type: %v", measurementType)
	} else {
		measurement.SetFinishTime()
	}

	return nil
}

func (reporter *Reporter) Measure(testCount *string, jsonName *string, measurementType MeasurementType, function func() (any, error)) (any, error) {
	reporter.StartMeasure(testCount, jsonName, measurementType)
	functionResult, functionError := function()
	if err := reporter.FinishMeasure(testCount, jsonName, measurementType); err != nil {
		panic(fmt.Sprintf("Finished measureting before started: testCount=%s ; jsonName=%s ; measurementType=%v", *testCount, *jsonName, measurementType))
	}
	return functionResult, functionError
}

func (reporter *Reporter) GetMeasures() ReportData {
	// Shaked-TODO: lock & clone the map before returning it
	return reporter.measurementDuration
}

func setDefault[K comparable, V any](mappy map[K]V, key K, defaultValue V) V {
	if result, ok := mappy[key]; ok {
		return result
	} else {
		mappy[key] = defaultValue
		return defaultValue
	}
}
