package Vector

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Vector[T any] struct {
	data []T
}

// #region Constructors

func NewVector[T any](length uint) *Vector[T] {
	return &Vector[T]{
		data: make([]T, length),
	}
}

func NewVector2[T any](capacity uint) *Vector[T] {
	return &Vector[T]{
		data: make([]T, 0, capacity),
	}
}

//goland:noinspection GoUnusedExportedFunction
func NewVector3[T any](length, capacity uint) *Vector[T] {
	return &Vector[T]{
		data: make([]T, length, capacity),
	}
}

// #endregion

// #region Setters

func (vector *Vector[T]) SetSlice(data []T) {
	vector.data = data
}

func (vector *Vector[T]) Push(element T) {
	vector.data = append(vector.data, element)
}

func (vector *Vector[T]) Pop() T {
	var element T
	element, vector.data = vector.data[len(vector.data)-1], vector.data[:len(vector.data)-1]
	return element
}

// #endregion

// #region Getters

func (vector *Vector[T]) GetAll() []T {
	return vector.data[:]
}

func (vector *Vector[T]) Len() int {
	return len(vector.data)
}

// #endregion

// #region (Un)Marshaler Interface

func (vector *Vector[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(vector.data)
}

// #endregion
