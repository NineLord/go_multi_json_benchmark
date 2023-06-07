package main

import (
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Config"
	"os"
)

// Example: make clean foo && clear  && ./bin/foo

func main() {
	println("Hello!")

	setRawConfig()

	println("Goodbye!")
}

func setRawConfig() {
	//config := Config.Config{
	//	Name:             "Shimi",
	//	Size:             "32KiB",
	//	Path:             "/mnt/c/Users/Shaked/Documents/Mine/IdeaProjects/rust_json_benchmark/junk/smallJson_2_n8d10m3.json",
	//	NumberOfLetters:  32,
	//	Depth:            8,
	//	NumberOfChildren: 3,
	//	Raw:              nil,
	//}
	configs := []Config.Config{
		{
			Name:             "Shimi",
			Size:             "32KiB",
			Path:             "/mnt/c/Users/Shaked/Documents/Mine/IdeaProjects/rust_json_benchmark/junk/smallJson_2_n8d10m3.json",
			NumberOfLetters:  32,
			Depth:            8,
			NumberOfChildren: 3,
			Raw:              nil,
		},
	}

	for index := range configs {
		config := &configs[index]
		var buffer []byte
		var err error
		if buffer, err = os.ReadFile(config.Path); err != nil {
			panic(err)
		}
		config.Raw = buffer
	}

	//println("Config: %+v", config)
	println("i hate go")
}

//goland:noinspection GoUnusedFunction
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

//goland:noinspection GoUnusedFunction
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
