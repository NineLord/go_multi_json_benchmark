package Reporter

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
