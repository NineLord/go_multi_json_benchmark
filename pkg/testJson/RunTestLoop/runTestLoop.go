package RunTestLoop

import (
	"fmt"
	"github.com/NineLord/go_multi_json_benchmark/pkg/searchTree/BreadthFirstSearch"
	"github.com/NineLord/go_multi_json_benchmark/pkg/searchTree/DepthFirstSearch"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter/MeasurementType"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter/StaticReporter"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/JsonGenerator"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var CharacterPoll = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz!@#$%&"

type RunTestLoop struct {
	testCount     uint
	valueToSearch float64
}

func NewRunTestLoop(testCount uint, valueToSearch float64) *RunTestLoop {
	return &RunTestLoop{
		testCount:     testCount,
		valueToSearch: valueToSearch,
	}
}

func (runTestLoop *RunTestLoop) RunTest(jsonName string, numberOfLetters uint, depth uint, numberOfChildren uint, rawJson []byte) error {
	for counter := uint(1); counter <= runTestLoop.testCount; counter++ {
		testName := fmt.Sprintf("Test %d", counter)
		if err := runTestLoop.runSingleTest(testName, jsonName, numberOfLetters, depth, numberOfChildren, rawJson); err != nil {
			return err
		}
	}
	return nil
}

func (runTestLoop *RunTestLoop) runSingleTest(testCount string, jsonName string, numberOfLetters uint, depth uint, numberOfChildren uint, rawJson []byte) error {
	return StaticReporter.Measure2(testCount, jsonName, MeasurementType.TotalIncludeContextSwitch, func() error {
		return runTestLoop.runSingleTestWithoutTotalMeasure(testCount, jsonName, numberOfLetters, depth, numberOfChildren, rawJson)
	})
}

func (runTestLoop *RunTestLoop) runSingleTestWithoutTotalMeasure(testCount string, jsonName string, numberOfLetters uint, depth uint, numberOfChildren uint, rawJson []byte) error {
	if _, err := testGeneratingJson(testCount, jsonName, numberOfLetters, depth, numberOfChildren); err != nil {
		return err
	} else if inputJsonFile, err := testDeserializeJson(testCount, jsonName, rawJson); err != nil {
		return err
	} else if err := runTestLoop.testIterateIteratively(testCount, jsonName, inputJsonFile); err != nil {
		return err
	} else if err := runTestLoop.testIterateRecursively(testCount, jsonName, inputJsonFile); err != nil {
		return err
	} else if _, err := testSerializeJson(testCount, jsonName, inputJsonFile); err != nil {
		return err
	}
	return nil
}

func testGeneratingJson(testCount string, jsonName string, numberOfLetters uint, depth uint, numberOfChildren uint) (map[string]interface{}, error) {
	await := make(chan utils.Result[map[string]interface{}])
	go func() {
		res, err := StaticReporter.Measure(testCount, jsonName, MeasurementType.GeneratingJson, func() (map[string]interface{}, error) {
			return JsonGenerator.GenerateJson(CharacterPoll, numberOfLetters, depth, numberOfChildren)
		})

		await <- utils.Result[map[string]interface{}]{
			Ok:  res,
			Err: err,
		}
	}()
	result := <-await
	return result.Ok, result.Err
}

func testDeserializeJson(testCount string, jsonName string, rawJson []byte) (map[string]interface{}, error) {
	await := make(chan utils.Result[map[string]interface{}])
	go func() {
		res, err := StaticReporter.Measure(testCount, jsonName, MeasurementType.DeserializeJson, func() (map[string]interface{}, error) {
			var inputJsonFile map[string]interface{}
			err := json.Unmarshal(rawJson, &inputJsonFile)
			return inputJsonFile, err
		})

		await <- utils.Result[map[string]interface{}]{
			Ok:  res,
			Err: err,
		}
	}()
	result := <-await
	return result.Ok, result.Err
}

func (runTestLoop *RunTestLoop) testIterateIteratively(testCount string, jsonName string, inputJsonFile map[string]interface{}) error {
	await := make(chan error)
	go func() {
		await <- StaticReporter.Measure2(testCount, jsonName, MeasurementType.IterateIteratively, func() error {
			if BreadthFirstSearch.Run(inputJsonFile, runTestLoop.valueToSearch) {
				return fmt.Errorf("BFS the tree found value that shouldn't be in it: %f", runTestLoop.valueToSearch)
			}
			return nil
		})
	}()
	return <-await
}

func (runTestLoop *RunTestLoop) testIterateRecursively(testCount string, jsonName string, inputJsonFile map[string]interface{}) error {
	await := make(chan error)
	go func() {
		await <- StaticReporter.Measure2(testCount, jsonName, MeasurementType.IterateRecursively, func() error {
			if DepthFirstSearch.Run(inputJsonFile, runTestLoop.valueToSearch) {
				return fmt.Errorf("DFS the tree found value that shouldn't be in it: %f", runTestLoop.valueToSearch)
			}
			return nil
		})
	}()
	return <-await
}

func testSerializeJson(testCount string, jsonName string, inputJsonFile map[string]interface{}) (string, error) {
	await := make(chan utils.Result[string])
	go func() {
		res, err := StaticReporter.Measure(testCount, jsonName, MeasurementType.SerializeJson, func() (string, error) {
			var buff []byte
			var err error
			if buff, err = json.Marshal(inputJsonFile); err != nil {
				return "", err
			}
			return string(buff), nil
		})

		await <- utils.Result[string]{
			Ok:  res,
			Err: err,
		}
	}()
	result := <-await
	return result.Ok, result.Err
}
