package main

import (
	// "encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"time"

	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/ExcelGenerator"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/PcUsageExporter"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/MathDataCollector"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/Vector"
	jsoniter "github.com/json-iterator/go"
	"github.com/struCoder/pidusage"
	"github.com/xuri/excelize/v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Example: make clean foo && clear  && ./bin/foo

func main() {
	divide(1,2)
}

func divide(dividend int, divisor int) (float64, error) {
	if divisor == 0 {
		return math.NaN(), errors.New("trying to divide by 0")
	} else {
		return float64(dividend) / float64(divisor), nil
	}
}


// func test17HighLevelProgramming(x ...int) (...[]interface{}) {
// 	return 0, 0
// }

func test17Oop() {
	myVec1 := Vector.NewVector[int](0)
	innerSlice1 := myVec1.GetAll()
	fmt.Println("innerSlice1", len(innerSlice1), cap(innerSlice1))
	myVec1.Push(1)
	fmt.Println(myVec1)
	myVec2 := Vector.Vector[int]{}
	innerSlice2 := myVec2.GetAll()
	fmt.Println("innerSlice2", len(innerSlice2), cap(innerSlice2))
	myVec2.Push(1)
	fmt.Println(myVec2)
}

func test16ExcelGenerator() {
	excelGenerator, err := ExcelGenerator.NewExcelGenerator("/tmp/some.json", 10, 3, 2, 2)
	if err != nil {
		panic(err)
	}
	measures := make(map[string]int64)
	measures["Test Generating JSON"] = 1
	measures["Test Deserialize JSON"] = 2
	measures["Test Iterate Iteratively"] = 3
	measures["Test Iterate Recursively"] = 4
	measures["Test Serialize JSON"] = 5
	pcUsages := []PcUsageExporter.PcUsage{
		{Cpu: 10.0, Ram: 1000},
		{Cpu: 25.0, Ram: 1500},
		{Cpu: 50.0, Ram: 2000},
	}
	if err := excelGenerator.AppendWorksheet("Testing 1", measures, pcUsages); err != nil {
		panic(err)
	}
	measures = make(map[string]int64)
	measures["Test Generating JSON"] = 2
	measures["Test Deserialize JSON"] = 3
	measures["Test Iterate Iteratively"] = 4
	measures["Test Iterate Recursively"] = 5
	measures["Test Serialize JSON"] = 6
	pcUsages = []PcUsageExporter.PcUsage{
		{Cpu: 10.0, Ram: 200},
		{Cpu: 25.0, Ram: 300},
		{Cpu: 25.0, Ram: 400},
	}
	if err := excelGenerator.AppendWorksheet("Testing 2", measures, pcUsages); err != nil {
		panic(err)
	}
	if err := excelGenerator.SaveAs("/mnt/c/Users/Shaked/Documents/Mine/IdeaProjects/go_json_benchmark/junk/report.xlsx"); err != nil {
		panic(err)
	}
}

func test15MathDataCollector() {
	dataCollector := MathDataCollector.MakeMathDataCollector()
	dataCollector.Add(1.0)
	fmt.Println(dataCollector.Minimum)
}

func test14PlayWithExcel() {
	report := excelize.NewFile()
	defer func() {
		if err := report.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := report.NewSheet("TestingSheetName"); err != nil {
		panic(err)
	}
	if err := report.SetPanes("TestingSheetName", &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	}); err != nil {
		panic(err)
	}
	if err := report.DeleteSheet("Sheet1"); err != nil {
		panic(err)
	}
	if err := report.SetCellDefault("TestingSheetName", "B4", "hello"); err != nil {
		panic(err)
	}
	if err := report.SetRowHeight("TestingSheetName", 4, 100); err != nil {
		panic(err)
	}
	if style, err := report.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	}); err != nil {
		panic(err)
	} else if err := report.SetCellStyle("TestingSheetName", "B4", "B4", style); err != nil {
		panic(err)
	}
	if err := report.SaveAs("/mnt/c/Users/Shaked/Documents/Mine/IdeaProjects/go_json_benchmark/junk/report.xlsx"); err != nil {
		panic(err)
	}
}

func test13ThreadGetMeTimeForFkSake() {
	receiver := make(chan []int64)
	sender := make(chan bool)
	go func(sender chan []int64, receiver chan bool, interval uint) {
		durationInterval := time.Duration(interval)
		result := make([]int64, 0)
		for {
			select {
			case signal := <-receiver:
				fmt.Printf("Thread: received a closing signal! signal: %t\n", signal)
				sender <- result
				close(sender)
				return
			default:
				result = append(result, time.Now().UnixMilli())
				time.Sleep(durationInterval * time.Millisecond)
			}
		}
	}(receiver, sender, 100)
	time.Sleep(500 * time.Millisecond)
	sender <- true
	close(sender)
	result := <-receiver
	for _, res := range result {
		fmt.Println(res)
	}
}

func test12ReadChannel() {
	fmt.Println("Test start!")
	myChannel := make(chan int)
	myChannel <- 1
	myChannel <- 2
	myChannel <- 3
	fmt.Println("Done inserting!")
	for index := 0; index < 3; index++ {
		fmt.Println(<-myChannel)
	}
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			close(c)
			fmt.Println("quit")
			return
		}
	}
}

func test11Fib() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for msg := range c {
			fmt.Println(msg)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

func test10ThreadGetTime() {
	receiver := make(chan bool)
	sender := make(chan bool)
	fmt.Println("Channels created!")
	go getTime(receiver, sender, 100)
	time.Sleep(1500 * time.Millisecond)
	fmt.Println("Closing thread")
	sender <- true
	close(sender)
	fmt.Println("Starting to collect messages")
	for message := range receiver {
		fmt.Printf("Message: %t\n", message)
	}
	fmt.Println("Main thread has finished!")
}

func getTime(sender chan bool, receiver chan bool, interval uint) {
	// timeToSend := time.Now().UnixMilli()
	durationInterval := time.Duration(interval)
	for {
		fmt.Println("Thread: First line of the loop")
		select {
		// case sender <- timeToSend:
		// 	fmt.Println("Thread: Sent the time!")
		// 	time.Sleep(durationInterval * time.Millisecond)
		// 	timeToSend = time.Now().UnixMilli()
		case signal := <-receiver:
			fmt.Printf("Thread: received a closing signal! signal: %t\n", signal)
			// close(sender)
			return
		default:
			sender <- false
			fmt.Println("Thread: Sent the time!")
			time.Sleep(durationInterval * time.Millisecond)
			// timeToSend = time.Now().UnixMilli()
			// fmt.Println("Thread: HELP ME")
		}

	}
}

func test9ThreadSleep() {
	go printTime(100)
	time.Sleep(500 * time.Millisecond)
}

func printTime(interval uint) {
	durationInterval := time.Duration(interval)
	for {
		fmt.Println(time.Now().UnixMilli())
		time.Sleep(durationInterval * time.Millisecond)
	}
}

func test8PcUsage() {
	sysInfo, err := pidusage.GetStat(os.Getpid())
	if err != nil {
		panic(err)
	}
	fmt.Printf("CPU: %.3f%% \t RAM: %.3fMB\n", sysInfo.CPU, (sysInfo.Memory/1024)/1024)
	fmt.Println(sysInfo)
}

func test7TypeOfNull() {
	fmt.Println(reflect.TypeOf([]int{1, 2, 3}).Kind())
	// fmt.Println(reflect.TypeOf(nil).Kind()) // Panic cus gO sMaRt!
	x := reflect.TypeOf(nil)
	fmt.Println(x)
}

func test6WhatTypeIsVector() {
	vec := Vector.NewVector[int](0)
	vec.Push(1)
	vec.Push(2)

	myType := reflect.TypeOf(vec)
	fmt.Println(myType)
}

func test5Marshal() {
	arr := []int{1, 2, 3}
	res, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	println(string(res))
}

func test4Stringify() {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(`{"a": 1, "b": {"c": 2, "d": {"e": 3}}}`), &data); err != nil {
		panic(err)
	}
	// jsonAsString, err := json.MarshalIndent(data, "", "  ")
	jsonAsString, err := json.MarshalToString(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonAsString))
}

func test3AddingToArray() {
	arrayies := make([]int, 0)
	arrayies = addItSlice(arrayies)
	arrayies = addItSlice(arrayies)
	arrayies = addItSlice(arrayies)
	arrayies = addItSlice(arrayies)
	arrayies = addItSlice(arrayies)
	arrayies = addItSlice(arrayies)
	fmt.Println(arrayies, len(arrayies), cap(arrayies))
}

func addItSlice(mySlice []int) []int {
	return append(mySlice, 5)
}

func test2AddingToMap() {
	mapy := make(map[string]interface{})
	addIt(mapy)
	fmt.Println(mapy)
}

func addIt(myMap map[string]interface{}) {
	myMap["a"] = 5
}

func test1() {
	// x := make([]int, 0, 3)
	z := [3]int{0, 0, 0}
	x := z[:0]
	y := append(x, 1)
	y = append(x, 2)
	fmt.Println("z:", z)
	innerMain(x)
	fmt.Println(y)
}

func innerMain(x interface{}) {
	switch reflect.TypeOf(x).Kind() {
	case reflect.Array:
		doSomethingWithArray(x.([3]int))
	case reflect.Slice:
		doSomethingWithSlice(x.([]int))
	default:
		panic("Wrong type!")
	}
}

func doSomethingWithSlice(arr []int) {
	fmt.Println("Slice", arr, len(arr), cap(arr))
}

func doSomethingWithArray(arr [3]int) {
	fmt.Println("Array", arr)
}
