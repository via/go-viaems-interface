package viaems

import (
	"time"
)

type StatusTarget interface {
	GetStatusUpdates() chan Status
}

type StatusLog interface {
	LatestStatus() *Status
	RangeOfStatuses(start time.Time, stop time.Time) []*Status
}

type SensorStatus struct {
	Value float64
	Fault bool
}

type Status struct {
	Sensors  map[string]*SensorStatus
	Fueling  map[string]float64
	Decoder  map[string]string
	Ignition map[string]float64
	CpuTime  float64
	WallTime time.Time
}

type TargetLogFile struct {
	path string
}
