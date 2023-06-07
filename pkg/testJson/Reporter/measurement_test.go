package Reporter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMeasurement_SetFinishTime(t *testing.T) {
	assertions := assert.New(t)

	measurement := MakeMeasurement()
	time.Sleep(1 * time.Second)
	measurement.SetFinishTime()
	durationPtr := measurement.GetDuration()
	assertions.NotNil(durationPtr)
	duration := *durationPtr

	assertions.GreaterOrEqual(duration, 1*time.Second)
}

func TestMeasurement_GetDuration(t *testing.T) {
	assertions := assert.New(t)

	measurement := MakeMeasurement()
	duration := measurement.GetDuration()

	assertions.Nil(duration)
}

func TestMeasurement_New_SetFinishTime(t *testing.T) {
	assertions := assert.New(t)

	measurement := NewMeasurement()
	time.Sleep(1 * time.Second)
	measurement.SetFinishTime()
	durationPtr := measurement.GetDuration()
	assertions.NotNil(durationPtr)
	duration := *durationPtr

	assertions.GreaterOrEqual(duration, 1*time.Second)
}

func TestMeasurement_New_GetDuration(t *testing.T) {
	assertions := assert.New(t)

	measurement := NewMeasurement()
	duration := measurement.GetDuration()

	assertions.Nil(duration)
}
