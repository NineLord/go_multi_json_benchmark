package main

// #region Imports
import (
	"errors"
	"fmt"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Config"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// #endregion

var CharacterPoll = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz!@#$%&"

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
	args := arguments.Args()
	var configs []Config.Config
	var testCounter uint
	if configs, testCounter, err = parseArgument(&args); err != nil {
		return
	}
	pathToSaveFile := arguments.Path("saveFile")
	threadCount := arguments.Uint("threadCount")
	debug := arguments.Bool("debug")

	runtime.GOMAXPROCS(int(threadCount))

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

	//var excelGenerator *ExcelGenerator.ExcelGenerator
	//if excelGenerator, err = ExcelGenerator.NewExcelGenerator(jsonPath, sampleInterval, numberOfLetters, depth, numberOfChildren); err != nil {
	//	return
	//}
	//var buffer []byte
	//if buffer, err = os.ReadFile(jsonPath); err != nil {
	//	return
	//}
	//valueToSearch := float64(2_000_000_000)
	//
	//for count := uint(0); count < testCounter; count++ {
	//	// #region Test Preparations
	//	reporter := Reporter.NewReporter()
	//	mainSender := make(chan bool)
	//	threadSender := make(chan []PcUsageExporter.PcUsage)
	//	go PcUsageExporter.Main(threadSender, mainSender, sampleInterval)
	//
	//	// #endregion
	//
	//	// #region Testing
	//
	//	if _, err = reporter.Measure("Test Generating JSON", func() (any, error) {
	//		return JsonGenerator.GenerateJson(CharacterPoll, numberOfLetters, depth, numberOfChildren)
	//	}); err != nil {
	//		return
	//	}
	//
	//	var inputJsonFile map[string]interface{}
	//	if _, err = reporter.Measure("Test Deserialize JSON", func() (any, error) {
	//		return nil, json.Unmarshal(buffer, &inputJsonFile)
	//	}); err != nil {
	//		return
	//	}
	//
	//	if found, _ := reporter.Measure("Test Iterate Iteratively", func() (any, error) {
	//		return BreadthFirstSearch.Run(inputJsonFile, valueToSearch), nil
	//	}); found == true {
	//		return fmt.Errorf("BFS the tree found value that shouldn't be in it: %f", valueToSearch)
	//	}
	//
	//	if found, _ := reporter.Measure("Test Iterate Recursively", func() (any, error) {
	//		return DepthFirstSearch.Run(inputJsonFile, valueToSearch), nil
	//	}); found == true {
	//		return fmt.Errorf("DFS the tree found value that shouldn't be in it: %f", valueToSearch)
	//	}
	//
	//	if _, err = reporter.Measure("Test Serialize JSON", func() (any, error) {
	//		var buff []byte
	//		var err error
	//		if buff, err = json.Marshal(inputJsonFile); err != nil {
	//			return nil, err
	//		}
	//		return string(buff), nil
	//	}); err != nil {
	//		return
	//	}
	//
	//	mainSender <- true
	//	close(mainSender)
	//	pcUsages := <-threadSender
	//	/*for _, pcUsage := range pcUsages {
	//		fmt.Printf("CPU: %.3f%% \t RAM: %.3fMB\n", pcUsage.Cpu, pcUsage.Ram)
	//	}
	//	for measureName, measureDuration := range reporter.GetMeasures() {
	//		fmt.Printf("%s : \t %d\n", measureName, measureDuration)
	//	}*/
	//	if err = excelGenerator.AppendWorksheet("Test "+strconv.Itoa(int(count)+1), reporter.GetMeasures(), pcUsages); err != nil {
	//		return
	//	}
	//
	//	// #endregion
	//}
	//
	//err = excelGenerator.SaveAs(pathToSaveFile)
	return
}
