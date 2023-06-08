package ExcelGenerator

import (
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Config"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter/MeasurementType"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/MathDataCollector"
	"github.com/xuri/excelize/v2"
	"time"
)

type ExcelGenerator struct {
	aboutInformation   []Config.Config
	workbook           *excelize.File
	formatBorder       int
	formatBorderCenter int
	formatColorful     int
	jsonNames          *utils.Vector[string]
	totalTestLength    time.Duration
	averagePerJsons    map[string]map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector
	averageAllJsons    map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector
}

func MakeExcelGenerator(jsonNames *utils.Vector[string], totalTestLength time.Duration, configs []Config.Config) (ExcelGenerator, error) {
	workbook := excelize.NewFile()

	formatBorder, err := workbook.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return ExcelGenerator{}, err
	}
	formatBorderCenter, err := workbook.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return ExcelGenerator{}, err
	}
	formatColorful, err := workbook.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"9AA9F6"}, Pattern: 1},
	})
	if err != nil {
		return ExcelGenerator{}, err
	}

	averagePerJsons := make(map[string]map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector)
	for _, jsonName := range jsonNames.GetAll() {
		averagePerJsons[jsonName] = getDataCollectorsForEachTest()
	}

	return ExcelGenerator{
		aboutInformation:   configs,
		workbook:           workbook,
		formatBorder:       formatBorder,
		formatBorderCenter: formatBorderCenter,
		formatColorful:     formatColorful,
		jsonNames:          jsonNames,
		totalTestLength:    totalTestLength,
		averagePerJsons:    averagePerJsons,
		averageAllJsons:    getDataCollectorsForEachTest(),
	}, nil
}

func NewExcelGenerator(jsonNames *utils.Vector[string], totalTestLength time.Duration, configs []Config.Config) (*ExcelGenerator, error) {
	if result, err := MakeExcelGenerator(jsonNames, totalTestLength, configs); err != nil {
		return nil, err
	} else {
		return &result, nil
	}
}

// // #region Adding Data
//
//	func (excelGenerator *ExcelGenerator) AppendWorksheet(worksheetName string, measures map[string]int64, pcUsage []PcUsageExporter.PcUsage) error {
//		if _, err := excelGenerator.workbook.NewSheet(worksheetName); err != nil {
//			return err
//		}
//		if err := excelGenerator.workbook.SetPanes(worksheetName, &excelize.Panes{
//			Freeze:      true,
//			YSplit:      1,
//			TopLeftCell: "A2",
//			ActivePane:  "bottomLeft",
//		}); err != nil {
//			return err
//		}
//
//		if err := excelGenerator.generateTitles(worksheetName); err != nil {
//			return err
//		}
//		if err := excelGenerator.addData(worksheetName, measures, pcUsage); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func (excelGenerator *ExcelGenerator) generateTitles(worksheetName string) error {
//		// #region Column 1
//
//		// #region Table 1
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A1", "Title", excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A2", "Generating JSON", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A3", "Iterating JSON Iteratively - BFS", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A4", "Iterating JSON Recursively - DFS", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A5", "Deserializing JSON", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A6", "Serializing JSON", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A7", "Total", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		// #region Table 2
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A9", "Average CPU (%)", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A10", "Average RAN (MB)", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		// #endregion
//
//		// #region Column 2
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "B1", "Time (ms)", excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		// #region Column 4
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "D1", "CPU (%)", excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		// #region Column 5
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "E1", "RAM (MB)", excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		return nil
//	}
//
//	func (excelGenerator *ExcelGenerator) addData(worksheetName string, measures map[string]int64, pcUsages []PcUsageExporter.PcUsage) error {
//		rowTotal := MathDataCollector.NewMathDataCollector()
//		columnCpuUsage := MathDataCollector.NewMathDataCollector()
//		columnRamUsage := MathDataCollector.NewMathDataCollector()
//
//		// #region JSON Manipulations
//
//		for testName, testResult := range measures {
//			intTestResult := int(testResult)
//			float64TestResult := float64(testResult)
//			switch testName {
//			case "Test Generating JSON":
//				if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B2", intTestResult, excelGenerator.formatBorderCenter); err != nil {
//					return err
//				}
//				excelGenerator.mathDataCollectors.averageGeneratingJsons.Add(float64TestResult)
//			case "Test Deserialize JSON":
//				if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B3", intTestResult, excelGenerator.formatBorderCenter); err != nil {
//					return err
//				}
//				excelGenerator.mathDataCollectors.averageIteratingJsonsIteratively.Add(float64TestResult)
//			case "Test Iterate Iteratively":
//				if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B4", intTestResult, excelGenerator.formatBorderCenter); err != nil {
//					return err
//				}
//				excelGenerator.mathDataCollectors.averageIteratingJsonsRecursively.Add(float64TestResult)
//			case "Test Iterate Recursively":
//				if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B5", intTestResult, excelGenerator.formatBorderCenter); err != nil {
//					return err
//				}
//				excelGenerator.mathDataCollectors.averageDeserializingJsons.Add(float64TestResult)
//			case "Test Serialize JSON":
//				if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B6", intTestResult, excelGenerator.formatBorderCenter); err != nil {
//					return err
//				}
//				excelGenerator.mathDataCollectors.averageSerializingJsons.Add(float64TestResult)
//			default:
//				return fmt.Errorf("invalid test name: %s", testName)
//			}
//			rowTotal.Add(float64TestResult)
//		}
//
//		if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B7", int(rowTotal.Sum), excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//
//		// #endreigon
//
//		// #region pc usage
//
//		currentRowNumber := 2
//		for _, pcUsage := range pcUsages {
//			if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "D"+strconv.Itoa(currentRowNumber), pcUsage.Cpu, excelGenerator.formatBorderCenter); err != nil {
//				return err
//			}
//			if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "E"+strconv.Itoa(currentRowNumber), pcUsage.Ram, excelGenerator.formatBorderCenter); err != nil {
//				return err
//			}
//
//			columnCpuUsage.Add(pcUsage.Cpu)
//			columnRamUsage.Add(pcUsage.Ram)
//
//			excelGenerator.mathDataCollectors.totalAverageCpu.Add(pcUsage.Cpu)
//			excelGenerator.mathDataCollectors.totalAverageRam.Add(pcUsage.Ram)
//
//			currentRowNumber++
//		}
//
//		if average, err := columnCpuUsage.Average(); err == nil {
//			if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "B9", average, excelGenerator.formatBorderCenter); err != nil {
//				return err
//			}
//		}
//		if average, err := columnRamUsage.Average(); err == nil {
//			if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "B10", average, excelGenerator.formatBorderCenter); err != nil {
//				return err
//			}
//		}
//
//		// #endregion
//
//		return nil
//	}
//
// // #endregion
//
// // #region Add Summary Worksheet
//
//	func (excelGenerator *ExcelGenerator) addAverageWorksheet() error {
//		worksheetName := "Average"
//		if _, err := excelGenerator.workbook.NewSheet(worksheetName); err != nil {
//			return err
//		}
//		if err := excelGenerator.generateAverageTitles(worksheetName); err != nil {
//			return err
//		}
//		return excelGenerator.addAverageData(worksheetName)
//	}
//
//	func (excelGenerator *ExcelGenerator) generateAverageTitles(worksheetName string) error {
//		// #region Column 1
//
//		// #region Table 1
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A1", "Title", excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A2", "Average Generating JSONs", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A3", "Average Iterating JSONs Iteratively - BFS", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A4", "Average Iterating JSONs Recursively - DFS", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A5", "Average Deserializing JSONs", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A6", "Average Serializing JSONs", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A7", "Average Totals", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		// #region Table 2
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A9", "Average Total CPU (%)", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A10", "Average Total RAN (MB)", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		// #endregion
//
//		// #region Column 2
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "B1", "Time (ms)", excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//
//		// #endregion
//
//		return nil
//	}
//
//	func (excelGenerator *ExcelGenerator) addAverageData(worksheetName string) error {
//		totalAverages := MathDataCollector.NewMathDataCollector()
//		cells := make(map[int]*MathDataCollector.MathDataCollector)
//		cells[2] = &excelGenerator.mathDataCollectors.averageGeneratingJsons
//		cells[3] = &excelGenerator.mathDataCollectors.averageIteratingJsonsIteratively
//		cells[4] = &excelGenerator.mathDataCollectors.averageIteratingJsonsRecursively
//		cells[5] = &excelGenerator.mathDataCollectors.averageDeserializingJsons
//		cells[6] = &excelGenerator.mathDataCollectors.averageSerializingJsons
//		for cellRow, dataCollector := range cells {
//			average, err := dataCollector.Average()
//			if err != nil {
//				continue
//			}
//			if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "B"+strconv.Itoa(cellRow), average, excelGenerator.formatBorderCenter); err != nil {
//				return nil
//			}
//			totalAverages.Add(average)
//		}
//		if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "B7", totalAverages.Sum, excelGenerator.formatBorderCenter); err != nil {
//			return err
//		}
//
//		if average, err := excelGenerator.mathDataCollectors.totalAverageCpu.Average(); err == nil {
//			if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "B9", average, excelGenerator.formatBorderCenter); err != nil {
//				return err
//			}
//		}
//
//		if average, err := excelGenerator.mathDataCollectors.totalAverageRam.Average(); err == nil {
//			if err := setCellFloat64AndStyle(excelGenerator.workbook, worksheetName, "B10", average, excelGenerator.formatBorderCenter); err != nil {
//				return err
//			}
//		}
//
//		return nil
//	}
//
// // #endregion
//
// // #region Add About Worksheet
//
//	func (excelGenerator *ExcelGenerator) addAboutWorksheet() error {
//		worksheetName := "About"
//		if _, err := excelGenerator.workbook.NewSheet(worksheetName); err != nil {
//			return err
//		}
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A1", "Path to JSON to be tested on (Iterating/Deserializing/Serializing)", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "B1", excelGenerator.aboutInformation.jsonPath, excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A2", "CPU/RAM Sampling Interval (milliseconds)", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B2", int(excelGenerator.aboutInformation.sampleInterval), excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A3", "Number of letters to generate for each node in the generated JSON tree", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B3", int(excelGenerator.aboutInformation.numberOfLetters), excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A4", "Depth of the generated JSON tree", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B4", int(excelGenerator.aboutInformation.depth), excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, "A5", "Number of children each node in the generated JSON tree going to have", excelGenerator.formatBorder); err != nil {
//			return err
//		}
//		if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, "B5", int(excelGenerator.aboutInformation.numberOfChildren), excelGenerator.formatBorder); err != nil {
//			return err
//		}
//
//		return nil
//	}
//
// // #endregion
//
// // #region Helper methods for excelize
//
//	func setCellDefaultAndStyle(workbook *excelize.File, worksheetName, cell, value string, style int) error {
//		if err := workbook.SetCellDefault(worksheetName, cell, value); err != nil {
//			return err
//		}
//		if err := workbook.SetCellStyle(worksheetName, cell, cell, style); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func setCellIntAndStyle(workbook *excelize.File, worksheetName, cell string, value, style int) error {
//		if err := workbook.SetCellInt(worksheetName, cell, value); err != nil {
//			return err
//		}
//		if err := workbook.SetCellStyle(worksheetName, cell, cell, style); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func setCellFloat64AndStyle(workbook *excelize.File, worksheetName, cell string, value float64, style int) error {
//		if err := workbook.SetCellFloat(worksheetName, cell, value, 2, 64); err != nil {
//			return err
//		}
//		if err := workbook.SetCellStyle(worksheetName, cell, cell, style); err != nil {
//			return err
//		}
//		return nil
//	}
//
// // #endregion
func (excelGenerator *ExcelGenerator) SaveAs(pathToSaveFile string) error {
	//if err := excelGenerator.addAverageWorksheet(); err != nil {
	//	return err
	//}
	//if err := excelGenerator.addAboutWorksheet(); err != nil {
	//	return err
	//}

	if err := excelGenerator.workbook.DeleteSheet("Sheet1"); err != nil { // For some reason doing it doesn't work in the constructor
		return err
	}

	if err := excelGenerator.workbook.SaveAs(pathToSaveFile); err != nil {
		return err
	}
	return excelGenerator.workbook.Close()
}

func getDataCollectorsForEachTest() map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector {
	dataCollectors := make(map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector)
	for _, measurementType := range MeasurementType.VariantsMeasurementTypes {
		dataCollectors[measurementType] = MathDataCollector.MakeMathDataCollector()
	}
	return dataCollectors
}
