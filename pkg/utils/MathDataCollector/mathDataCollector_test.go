package MathDataCollector

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMathDataCollectorCorrectValues(t *testing.T) {
	assertions := assert.New(t)

	dataCollector := MakeMathDataCollector()
	dataCollector.Add(1.0)
	dataCollector.Add(2.0)
	dataCollector.Add(3.0)
	dataCollector.Add(2.0)
	dataCollector.Add(1.0)
	dataCollector.Add(2.0)
	dataCollector.Add(3.0)
	dataCollector.Add(2.0)
	dataCollector.Add(1.0)

	assertions.Equal(dataCollector.Maximum, 3.0)
	assertions.Equal(dataCollector.Minimum, 1.0)
	assertions.Equal(dataCollector.Sum, 1.0+2.0+3.0+2.0+1.0+2.0+3.0+2.0+1.0)
	assertions.Equal(dataCollector.Count, uint64(9))
	if avg, err := dataCollector.Average(); err != nil {
		assertions.Fail("Failed to get average", err)
	} else {
		assertions.Equal(avg, (1.0+2.0+3.0+2.0+1.0+2.0+3.0+2.0+1.0)/9.0)
	}
}
