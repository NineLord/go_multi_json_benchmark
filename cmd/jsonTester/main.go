package main

// #region Imports
import (
	"fmt"
	"github.com/NineLord/go_multi_json_benchmark/pkg/searchTree/BreadthFirstSearch"
	"github.com/NineLord/go_multi_json_benchmark/pkg/searchTree/DepthFirstSearch"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/ExcelGenerator"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/PcUsageExporter"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/JsonGenerator"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// #endregion

var CharacterPoll = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz!@#$%&"

func getDefaultPathToSaveFile() cli.Path {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir + "/report4.xlsx"
	} else if executable, err := os.Executable(); err == nil {
		return filepath.Dir(executable) + "/report4.xlsx"
	} else {
		panic(fmt.Sprintf("Didn't get result and couldn't get default path, error: %s", err))
	}
}

// Example: make clean jsonTester && clear && ./bin/jsonTester -i 10 -n 8 -d 10 -m 5 -s ./junk/report4.xlsx -D /mnt/c/Users/Shaked/Documents/Mine/IdeaProjects/rust_json_benchmark/junk/hugeJson_numberOfLetters8_depth10_children5.json 5

func main() {
	app := &cli.App{
		Name:      "jsonTester",
		Usage:     "Tests JSON manipulations",
		ArgsUsage: "<jsonPath> [testCounter]", // Can't add description to arguments with this package
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "saveFile",
				Aliases: []string{"s"},
				Usage:   "Absolute path to save the excel report file to",
				Value:   getDefaultPathToSaveFile(),
			},
			&cli.UintFlag{
				Name:    "sampleInterval",
				Aliases: []string{"i"},
				Usage:   "The interval in which it will sample the CPU/RAM usage of the system while running the tests, units are in milliseconds",
				Value:   50,
			},
			&cli.UintFlag{
				Name:    "numberOfLetters",
				Aliases: []string{"n"},
				Usage:   "The total number of letters that each generated node name will have in the generated JSON tree",
				Value:   32,
			},
			&cli.UintFlag{
				Name:    "depth",
				Aliases: []string{"d"},
				Usage:   "The depth of the generated JSON tree",
				Value:   255,
			},
			&cli.UintFlag{
				Name:    "numberOfChildren",
				Aliases: []string{"m"},
				Usage:   "The number of children each node should have in the generated JSON tree",
				Value:   16,
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

func parseArgument(arguments *cli.Args) (cli.Path, uint) {
	jsonPath := (*arguments).Get(0)
	if jsonPath == "" {
		panic("First argument must be the absolute path to the JSON file that will be tested")
	}
	if testCounter := (*arguments).Get(1); testCounter == "" {
		return jsonPath, 5
	} else if testCounter, err := strconv.Atoi(testCounter); err != nil {
		panic(fmt.Sprintf("Second argument should be the number of times will run the tests: %s", err))
	} else {
		return jsonPath, uint(testCounter)
	}
}

func cliAction(arguments *cli.Context) (err error) {
	args := arguments.Args()
	jsonPath, testCounter := parseArgument(&args)
	pathToSaveFile := arguments.Path("saveFile")
	sampleInterval := arguments.Uint("sampleInterval")
	numberOfLetters := arguments.Uint("numberOfLetters")
	depth := arguments.Uint("depth")
	numberOfChildren := arguments.Uint("numberOfChildren")
	debug := arguments.Bool("debug")

	if debug {
		fmt.Println(
			"Arguments:\n",
			"jsonPath:", jsonPath, "\n",
			"testCounter:", testCounter, "\n",
			"pathToSaveFile:", pathToSaveFile, "\n",
			"sampleInterval:", sampleInterval, "\n",
			"numberOfLetters:", numberOfLetters, "\n",
			"depth:", depth, "\n",
			"numberOfChildren:", numberOfChildren, "\n",
			"debug:", debug,
		)
	}

	var excelGenerator *ExcelGenerator.ExcelGenerator
	if excelGenerator, err = ExcelGenerator.NewExcelGenerator(jsonPath, sampleInterval, numberOfLetters, depth, numberOfChildren); err != nil {
		return
	}
	var buffer []byte
	if buffer, err = os.ReadFile(jsonPath); err != nil {
		return
	}
	valueToSearch := float64(2_000_000_000)

	for count := uint(0); count < testCounter; count++ {
		// #region Test Preparations
		reporter := Reporter.NewReporter()
		mainSender := make(chan bool)
		threadSender := make(chan []PcUsageExporter.PcUsage)
		go PcUsageExporter.Main(threadSender, mainSender, sampleInterval)

		// #endregion

		// #region Testing

		if _, err = reporter.Measure("Test Generating JSON", func() (any, error) {
			return JsonGenerator.GenerateJson(CharacterPoll, numberOfLetters, depth, numberOfChildren)
		}); err != nil {
			return
		}

		var inputJsonFile map[string]interface{}
		if _, err = reporter.Measure("Test Deserialize JSON", func() (any, error) {
			return nil, json.Unmarshal(buffer, &inputJsonFile)
		}); err != nil {
			return
		}

		if found, _ := reporter.Measure("Test Iterate Iteratively", func() (any, error) {
			return BreadthFirstSearch.Run(inputJsonFile, valueToSearch), nil
		}); found == true {
			return fmt.Errorf("BFS the tree found value that shouldn't be in it: %f", valueToSearch)
		}

		if found, _ := reporter.Measure("Test Iterate Recursively", func() (any, error) {
			return DepthFirstSearch.Run(inputJsonFile, valueToSearch), nil
		}); found == true {
			return fmt.Errorf("DFS the tree found value that shouldn't be in it: %f", valueToSearch)
		}

		if _, err = reporter.Measure("Test Serialize JSON", func() (any, error) {
			var buff []byte
			var err error
			if buff, err = json.Marshal(inputJsonFile); err != nil {
				return nil, err
			}
			return string(buff), nil
		}); err != nil {
			return
		}

		mainSender <- true
		close(mainSender)
		pcUsages := <-threadSender
		/*for _, pcUsage := range pcUsages {
			fmt.Printf("CPU: %.3f%% \t RAM: %.3fMB\n", pcUsage.Cpu, pcUsage.Ram)
		}
		for measureName, measureDuration := range reporter.GetMeasures() {
			fmt.Printf("%s : \t %d\n", measureName, measureDuration)
		}*/
		if err = excelGenerator.AppendWorksheet("Test "+strconv.Itoa(int(count)+1), reporter.GetMeasures(), pcUsages); err != nil {
			return
		}

		// #endregion
	}

	err = excelGenerator.SaveAs(pathToSaveFile)
	return
}
