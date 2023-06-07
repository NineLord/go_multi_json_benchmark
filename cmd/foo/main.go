package main

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Example: make clean foo && clear  && ./bin/foo

func main() {
	println("Hello!")

	awaitFunc2()

	println("Goodbye!")
}

func awaitFunc2() {
	foo := make(chan bool)
	go func() {
		cantUsePointAsMapKey()
		foo <- true
	}()
	<-foo
	println("TRUE")
	go func() {
		cantUsePointAsMapKey()
		foo <- true
	}()
	<-foo
	println("TRUE")
}

func awaitFunc() {
	foo := make(chan bool)
	go func() {
		cantUsePointAsMapKey()
		foo <- true
	}()
	x := <-foo
	println(x)
	go func() {
		cantUsePointAsMapKey()
		foo <- true
	}()
	x2 := <-foo
	println(x2)
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
