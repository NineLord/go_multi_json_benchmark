package main

//#region Imports
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Config"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/ExcelGenerator"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/RunTestLoop"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

//#endregion

func getDefaultPathToSaveFile() cli.Path {
	reportFileName := "/report_go.xlsx"
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir + reportFileName
	} else if executable, err := os.Executable(); err == nil {
		return filepath.Dir(executable) + reportFileName
	} else {
		panic(fmt.Sprintf("Didn't get result and couldn't get default path, error: %s", err))
	}
}

// Example: Shaked-TODO

func main() {
	app := &cli.App{
		Name:      "jsonTester",
		Usage:     "Tests JSON manipulations",
		ArgsUsage: "<configPath> [testCounter]", // Can't add description to arguments with this package
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "saveFile",
				Aliases: []string{"s"},
				Usage:   "Absolute path to save the excel report file to",
				Value:   getDefaultPathToSaveFile(),
			},
			&cli.UintFlag{
				Name:    "threadCount",
				Aliases: []string{"c"},
				Usage:   "Number of threads to use to run the test",
				Value:   3,
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"D"},
				Usage:   "Prints additional debug information",
				Value:   false,
			},
		},
		Action: cliAction,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseArgument(arguments *cli.Args) (configFile []Config.Config, testCounter uint, err error) {
	if (*arguments).Len() != 2 {
		//goland:noinspection GoErrorStringFormat
		err = errors.New("Too many arguments passed (did you put flags after argument?)")
		return
	}
	configPath := (*arguments).Get(0)
	if configPath == "" {
		err = errors.New("first argument must be the absolute path to the JSON file that will be tested")
		return
	}
	var buffer []byte
	if buffer, err = os.ReadFile(configPath); err != nil {
		return
	}
	if err = json.Unmarshal(buffer, &configFile); err != nil {
		return
	}

	var testCounterInt int
	if testCounterStr := (*arguments).Get(1); testCounterStr == "" {
		testCounter = 5
	} else if testCounterInt, err = strconv.Atoi(testCounterStr); err != nil {
		err = fmt.Errorf("second argument should be the number of times will run the tests: %s", err)
	} else {
		testCounter = uint(testCounterInt)
	}
	return
}

func cliAction(arguments *cli.Context) (err error) {
	//#region Input Parameters
	args := arguments.Args()
	var configs []Config.Config
	var testCounter uint
	if configs, testCounter, err = parseArgument(&args); err != nil {
		return
	}
	pathToSaveFile := arguments.Path("saveFile")
	threadCount := arguments.Uint("threadCount")
	debug := arguments.Bool("debug")
	//#endregion

	runtime.GOMAXPROCS(int(threadCount)) // Shaked-TODO: make this optional

	if debug {
		var prettyConfig string
		prettyConfigBuffer, err := json.MarshalIndent(configs, "", "    ")
		if err != nil {
			prettyConfig = fmt.Sprintf("%#v", configs)
		} else {
			prettyConfig = string(prettyConfigBuffer)
		}
		fmt.Println(
			"Arguments:\n",
			"configPath:", prettyConfig, "\n",
			"testCounter:", testCounter, "\n",
			"pathToSaveFile:", pathToSaveFile, "\n",
			"threadCount:", threadCount, "\n",
			"debug:", debug,
		)
	}

	//#region Test preparations
	testNames := utils.NewVector2[string](uint(len(configs)))
	for index := range configs {
		config := &configs[index]
		var buffer []byte
		if buffer, err = os.ReadFile(config.Path); err != nil {
			return
		}
		config.Raw = buffer
		testNames.Push(config.Name)
	}
	valueToSearch := float64(2_000_000_000)
	testRunner := RunTestLoop.NewRunTestLoop(testCounter, valueToSearch)
	errGroup := new(errgroup.Group)
	//#endregion

	//#region Testing
	measurementTotalTestLength := Reporter.MakeMeasurement()
	for index := range configs {
		config := &configs[index]
		errGroup.Go(func() error {
			return testRunner.RunTest(config.Name, config.NumberOfLetters, config.Depth, config.NumberOfChildren, config.Raw)
		})
	}
	if err = errGroup.Wait(); err != nil {
		return
	}
	measurementTotalTestLength.SetFinishTime()
	var totalTestLength time.Duration
	if duration := measurementTotalTestLength.GetDuration(); duration != nil {
		totalTestLength = *duration
	}
	//#endregion

	measurement := Reporter.GetInstance().GetMeasures()

	if debug {
		buffer, err := json.MarshalIndent(measurement, "", "    ")
		if err != nil {
			println("Can't show measurements")
		} else {
			println(string(buffer))
		}
		println("Whole test:", totalTestLength.Milliseconds())
	}

	var excelGenerator *ExcelGenerator.ExcelGenerator
	if excelGenerator, err = ExcelGenerator.NewExcelGenerator(testNames, totalTestLength, configs); err != nil {
		return
	}
	for counter := uint(1); counter <= testCounter; counter++ {
		testName := fmt.Sprintf("Test %d", counter)
		testCase, ok := measurement[testName]
		if !ok {
			err = fmt.Errorf("report doesn't contain the test name: %s", testName)
		}
		err = excelGenerator.AppendWorksheet(testName, testCase)
		if err != nil {
			return
		}
	}
	err = excelGenerator.SaveAs(pathToSaveFile)
	return
}
