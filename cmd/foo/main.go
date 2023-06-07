package main

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Example: make clean foo && clear  && ./bin/foo

func main() {
	println("Hello world!")
}
