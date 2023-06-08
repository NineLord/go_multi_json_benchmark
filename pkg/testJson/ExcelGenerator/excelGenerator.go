package ExcelGenerator

import (
	"fmt"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Config"
	"github.com/NineLord/go_multi_json_benchmark/pkg/testJson/Reporter/MeasurementType"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils"
	"github.com/NineLord/go_multi_json_benchmark/pkg/utils/MathDataCollector"
	"github.com/xuri/excelize/v2"
	"math"
	"strconv"
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

// #region Adding Data
func (excelGenerator *ExcelGenerator) AppendWorksheet(worksheetName string, measures map[string]map[MeasurementType.MeasurementType]time.Duration) (err error) {
	if _, err = excelGenerator.workbook.NewSheet(worksheetName); err != nil {
		return
	}

	testDataCollectors := getDataCollectorsForEachTest()
	currentRow := uint(1)

	for _, jsonName := range excelGenerator.jsonNames.GetAll() {
		testData, ok := measures[jsonName]
		if !ok {
			return fmt.Errorf("given database doesn't have a the JSON name: %s", jsonName)
		}

		if currentRow, err = excelGenerator.setColorfulTitle(worksheetName, currentRow, 1, jsonName); err != nil {
			return
		}

		jsonDataCollector := MathDataCollector.NewMathDataCollector()

		if currentRow, err = excelGenerator.addTestData(currentRow, 1, MeasurementType.GeneratingJson, "Generating JSON", worksheetName, jsonName, testData, jsonDataCollector, testDataCollectors); err != nil {
			return
		}
		if currentRow, err = excelGenerator.addTestData(currentRow, 1, MeasurementType.IterateIteratively, "Iterating JSON Iteratively - BFS", worksheetName, jsonName, testData, jsonDataCollector, testDataCollectors); err != nil {
			return
		}
		if currentRow, err = excelGenerator.addTestData(currentRow, 1, MeasurementType.IterateRecursively, "Iterating JSON Recursively - DFS", worksheetName, jsonName, testData, jsonDataCollector, testDataCollectors); err != nil {
			return
		}
		if currentRow, err = excelGenerator.addTestData(currentRow, 1, MeasurementType.DeserializeJson, "Deserializing JSON", worksheetName, jsonName, testData, jsonDataCollector, testDataCollectors); err != nil {
			return
		}
		if currentRow, err = excelGenerator.addTestData(currentRow, 1, MeasurementType.SerializeJson, "Serializing JSON", worksheetName, jsonName, testData, jsonDataCollector, testDataCollectors); err != nil {
			return
		}
		if currentRow, err = excelGenerator.addTotalTestData(currentRow, 1, MeasurementType.Total, "Total", worksheetName, jsonName, jsonDataCollector, testDataCollectors); err != nil {
			return
		}
		if currentRow, err = excelGenerator.addTestData(currentRow, 1, MeasurementType.TotalIncludeContextSwitch, "Total Including Context Switch", worksheetName, jsonName, testData, jsonDataCollector, testDataCollectors); err != nil {
			return
		}

		currentRow++
	}

	return
}

func (excelGenerator *ExcelGenerator) addTestData(row uint, column uint,
	measurementType MeasurementType.MeasurementType, title string,
	worksheetName, jsonName string,
	testData map[MeasurementType.MeasurementType]time.Duration,
	jsonDataCollector *MathDataCollector.MathDataCollector,
	testDataCollectors map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector) (uint, error) {

	rowString := strconv.Itoa(int(row))
	cell := columnNumberToString(int(column)) + rowString
	nextCell := columnNumberToString(int(column)+1) + rowString

	var value float64
	if duration, ok := testData[measurementType]; !ok {
		return row, fmt.Errorf("given database doesn't have a measurement type: %d", measurementType)
	} else {
		value = float64(duration.Milliseconds())
	}
	if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, cell, title, excelGenerator.formatBorder); err != nil {
		return row, err
	}
	if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, nextCell, int(value), excelGenerator.formatBorderCenter); err != nil {
		return row, err
	}
	jsonDataCollector.Add(value)
	if measureMap, ok := excelGenerator.averagePerJsons[jsonName]; !ok {
		return row, fmt.Errorf("averages_per_jsons doesn't have the given JSON name: %s", jsonName)
	} else if dataCollector, ok := measureMap[measurementType]; !ok {
		return row, fmt.Errorf("averages_per_jsons doesn't have the given measurement type: %d", measurementType)
	} else {
		dataCollector.Add(value)
	}
	if dataCollector, ok := excelGenerator.averageAllJsons[measurementType]; !ok {
		return row, fmt.Errorf("averages_all_jsons doesn't have the given measurement type: %d", measurementType)
	} else {
		dataCollector.Add(value)
	}
	if dataCollector, ok := testDataCollectors[measurementType]; !ok {
		return row, fmt.Errorf("test_data_collectors doesn't have the given measurement type: %d", measurementType)
	} else {
		dataCollector.Add(value)
	}

	return row + 1, nil
}

func (excelGenerator *ExcelGenerator) addTotalTestData(row uint, column uint,
	measurementType MeasurementType.MeasurementType, title string,
	worksheetName, jsonName string,
	jsonDataCollector *MathDataCollector.MathDataCollector,
	testDataCollectors map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector) (uint, error) {

	rowString := strconv.Itoa(int(row))
	cell := columnNumberToString(int(column)) + rowString
	nextCell := columnNumberToString(int(column)+1) + rowString

	value := jsonDataCollector.GetSum()
	if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, cell, title, excelGenerator.formatBorder); err != nil {
		return row, err
	}
	if err := setCellIntAndStyle(excelGenerator.workbook, worksheetName, nextCell, int(value), excelGenerator.formatBorderCenter); err != nil {
		return row, err
	}
	if measureMap, ok := excelGenerator.averagePerJsons[jsonName]; !ok {
		return row, fmt.Errorf("averages_per_jsons doesn't have the given JSON name: %s", jsonName)
	} else if dataCollector, ok := measureMap[measurementType]; !ok {
		return row, fmt.Errorf("averages_per_jsons doesn't have the given measurement type: %d", measurementType)
	} else {
		dataCollector.Add(value)
	}
	if dataCollector, ok := excelGenerator.averageAllJsons[measurementType]; !ok {
		return row, fmt.Errorf("averages_all_jsons doesn't have the given measurement type: %d", measurementType)
	} else {
		dataCollector.Add(value)
	}
	if dataCollector, ok := testDataCollectors[measurementType]; !ok {
		return row, fmt.Errorf("test_data_collectors doesn't have the given measurement type: %d", measurementType)
	} else {
		dataCollector.Add(value)
	}

	return row + 1, nil
}

//#endregion

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

// #region Helper methods for excelize
func (excelGenerator *ExcelGenerator) setColorfulTitle(worksheetName string, row uint, column uint, title string) (uint, error) {
	rowString := strconv.Itoa(int(row))
	cell := columnNumberToString(int(column)) + rowString
	nextCell := columnNumberToString(int(column)+1) + rowString

	if err := setCellDefaultAndStyle(excelGenerator.workbook, worksheetName, cell, title, excelGenerator.formatColorful); err != nil {
		return row, err
	}
	if err := excelGenerator.workbook.MergeCell(worksheetName, cell, nextCell); err != nil {
		return row, err
	}

	return row + 1, nil
}

func setCellDefaultAndStyle(workbook *excelize.File, worksheetName, cell, value string, style int) error {
	if err := workbook.SetCellDefault(worksheetName, cell, value); err != nil {
		return err
	}
	if err := workbook.SetCellStyle(worksheetName, cell, cell, style); err != nil {
		return err
	}
	return nil
}

func setCellIntAndStyle(workbook *excelize.File, worksheetName, cell string, value, style int) error {
	if err := workbook.SetCellInt(worksheetName, cell, value); err != nil {
		return err
	}
	if err := workbook.SetCellStyle(worksheetName, cell, cell, style); err != nil {
		return err
	}
	return nil
}

func setCellFloat64AndStyle(workbook *excelize.File, worksheetName, cell string, value float64, style int) error {
	if err := workbook.SetCellFloat(worksheetName, cell, value, 2, 64); err != nil {
		return err
	}
	if err := workbook.SetCellStyle(worksheetName, cell, cell, style); err != nil {
		return err
	}
	return nil
}

func getDataCollectorsForEachTest() map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector {
	dataCollectors := make(map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector)
	for _, measurementType := range MeasurementType.VariantsMeasurementTypes {
		dataCollectors[measurementType] = MathDataCollector.MakeMathDataCollector()
	}
	return dataCollectors
}

func columnStringToNumber(column string) uint {
	letters := []rune(column)
	result := uint(0)
	for _, letter := range letters {
		result = uint(letter) - 64 + result*26
	}
	return result
}

func columnNumberToString(column int) string {
	column--
	ordA := int('A')
	ordZ := int('Z')
	length := ordZ - ordA + 1

	result := ""
	for 0 <= column {
		x := string(rune(column%length + ordA))
		result = x + result
		column = int(math.Floor(float64(column)/float64(length))) - 1
	}
	return result
}

//#endregion
