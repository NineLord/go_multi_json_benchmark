package Reporter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReporter_Measure(t *testing.T) {
	assertions := assert.New(t)
	testCase := "Test 1"
	jsonName := "Json 1"
	measurementType := GeneratingJson

	reporter := GetInstance()
	if _, err := reporter.Measure(&testCase, &jsonName, measurementType, func() (any, error) {
		time.Sleep(1 * time.Second)
		return nil, nil
	}); err != nil {
		panic(err)
	}

	measurements := reporter.GetMeasures()
	if testMap, ok := measurements[testCase]; !ok {
		assertions.Fail("No test map")
	} else if jsonMap, ok := testMap[jsonName]; !ok {
		assertions.Fail("No json map")
	} else if measurement, ok := jsonMap[measurementType]; !ok {
		assertions.Fail("No duration for measurement type")
	} else if duration := measurement.GetDuration(); duration == nil {
		assertions.Fail("Measurement haven't finished")
	} else {
		duration := *duration
		assertions.GreaterOrEqual(duration, 1*time.Second)
		assertions.Less(duration, 2*time.Second)
	}
}
