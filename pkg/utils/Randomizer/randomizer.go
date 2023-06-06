package Randomizer

import (
	"math/rand"
)

func GetRandomBoolean() bool {
	return rand.Intn(2) == 1
}

func GetRandomNumberInRangeInt(minimum int, maximum int) int {
	return rand.Intn(maximum-minimum) + minimum
}

func GetRandomNumberInRangeFloat64(minimum float64, maximum float64) float64 {
	return minimum + rand.Float64()*(maximum-minimum)
}

func GetRandomIndexFromArray[T interface{}](array []T) int {
	return GetRandomNumberInRangeInt(0, len(array))
}

func GetRandomValueFromArray[T interface{}](array []T) *T {
	return &array[GetRandomIndexFromArray(array)]
}
