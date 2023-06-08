package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVector_Pop(t *testing.T) {
	assertions := assert.New(t)

	vector := NewVector[int](0)
	vector.Push(1)
	vector.Push(2)
	vector.Push(3)

	assertions.Equal([]int{1, 2, 3}, vector.GetAll())
}

func TestVector_MarshalJSON(t *testing.T) {
	assertions := assert.New(t)

	vector := NewVector[int](0)
	vector.Push(1)
	vector.Push(2)
	vector.Push(3)

	result, err := vector.MarshalJSON()
	if err != nil {
		assertions.Fail("Error while marshalling", err)
	}
	assertions.Equal("[1,2,3]", string(result))
}

func TestVector_MarshalJSON2(t *testing.T) {
	assertions := assert.New(t)

	vector := NewVector[int](0)
	vector.Push(1)
	vector.Push(2)
	vector.Push(3)

	result, err := json.Marshal(vector)
	if err != nil {
		assertions.Fail("Error while marshalling", err)
	}
	assertions.Equal("[1,2,3]", string(result))
}
