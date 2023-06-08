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

// #region Adding Data
func (excelGenerator *ExcelGenerator) AppendWorksheet(worksheetName string, measures map[string]map[MeasurementType.MeasurementType]time.Duration) error {
	if _, err := excelGenerator.workbook.NewSheet(worksheetName); err != nil {
		return err
	}
	return nil
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

func getDataCollectorsForEachTest() map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector {
	dataCollectors := make(map[MeasurementType.MeasurementType]MathDataCollector.MathDataCollector)
	for _, measurementType := range MeasurementType.VariantsMeasurementTypes {
		dataCollectors[measurementType] = MathDataCollector.MakeMathDataCollector()
	}
	return dataCollectors
}
