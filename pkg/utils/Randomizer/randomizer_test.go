package Randomizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRandomNumberInRange(t *testing.T) {
	assertions := assert.New(t)

	minimum, maximum := 0, 100
	for count := 0; count < 100; count++ {
		result := GetRandomNumberInRangeInt(minimum, maximum)
		assertions.GreaterOrEqual(result, minimum)
		assertions.Less(result, maximum)
	}
}

func TestGetRandomIndexFromArray(t *testing.T) {
	assertions := assert.New(t)

	array := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for count := 0; count < 100; count++ {
		result := GetRandomIndexFromArray(array)
		assertions.GreaterOrEqual(result, 0)
		assertions.Less(result, len(array))
	}
}
