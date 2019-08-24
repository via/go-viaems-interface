package main

import (
  "time"
)

type SensorStatus {
  Value float64
  Fault bool
}

type Status struct {
  Sensors map[string]SensorStatus
  Fueling map[string]float64
  Decoder map[string]float64
  Ignition map[string]float64
  CpuTime float64
  WallTime time.Time
}

type TargetLogFile struct {
  string path
}

