package Reporter

import "time"

type Measurement struct {
	startTime int64
	duration  *time.Duration
}

func MakeMeasurement() Measurement {
	return Measurement{
		startTime: time.Now().UnixMilli(),
		duration:  nil,
	}
}

func NewMeasurement() *Measurement {
	result := MakeMeasurement()
	return &result
}

func (measurement *Measurement) SetFinishTime() {
	result := time.Duration(time.Now().UnixMilli()-measurement.startTime) * time.Millisecond
	measurement.duration = &result
}

func (measurement *Measurement) GetDuration() *time.Duration {
	return measurement.duration
}
