package MeasurementType

type MeasurementType uint

const (
	GeneratingJson MeasurementType = iota
	DeserializeJson
	IterateIteratively
	IterateRecursively
	SerializeJson
	Total
	TotalIncludeContextSwitch
)

var VariantsMeasurementTypes = [7]MeasurementType{
	GeneratingJson,
	DeserializeJson,
	IterateIteratively,
	IterateRecursively,
	SerializeJson,
	Total,
	TotalIncludeContextSwitch,
}
