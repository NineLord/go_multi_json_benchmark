package main

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Example: make clean foo && clear  && ./bin/foo

func main() {
	println("Hello!")

	cantUsePointAsMapKey()

	println("Goodbye!")
}

func cantUsePointAsMapKey() {
	key1 := "a"
	key2 := "b"
	mapy := make(map[*string]int)
	mapy[&key1] = 1
	mapy[&key2] = 2
	key1 = "c"
	println("map1")
	key3 := "b"
	mapy[&key3] = 3
	println("map2")
}
