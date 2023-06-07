package Reporter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMeasurement_SetFinishTime(t *testing.T) {
	assertions := assert.New(t)

	measurement := NewMeasurement()
	time.Sleep(1 * time.Second)
	measurement.SetFinishTime()
	duration := *measurement.GetDuration()

	assertions.GreaterOrEqual(duration, 1*time.Second)
}

func TestMeasurement_GetDuration(t *testing.T) {
	assertions := assert.New(t)

	measurement := NewMeasurement()
	duration := measurement.GetDuration()

	assertions.Nil(duration)
}

func TestMakeMeasurement(t *testing.T) {
	assertions := assert.New(t)

	measurement := MakeMeasurement()
	time.Sleep(1 * time.Second)
	measurement.SetFinishTime()
	duration := *measurement.GetDuration()

	assertions.GreaterOrEqual(duration, 1*time.Second)
}
