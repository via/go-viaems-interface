package main

type ConfigurationTarget interface {
  ListTables() ([]string, error)
  GetTable(name string) (TableConfig, error)
}

type TableConfig struct {
  Name string
  AxisCount int
  RowCount int
  ColumnCount int

  ColumnLabels []string
  RowLabels []string

  ColumnName string
  RowName string
}

type EventConfig struct {
}

type SensorConfig struct {
}

type FuelingConfig struct {
}

type DecoderConfig struct {
}

type IgnitionConfig struct {
}

type TargetConfig struct {
  Tables map[string]TableConfig
  Events []EventConfig
  Sensors map[string]SensorConfig

  fueling FuelingConfig
  ignition IgnitionConfig
  decoder DecoderConfig
}

func NewTargetConfigFromFile(path string) *TargetConfig {
  return nil
}

func NewTargetConfigFromWire(w *ConfigurationTarget) *TargetConfig {
  return nil
}

func (t *TargetConfig) WriteTargetConfigToFile(path string) error {
  return nil
}

func (t *TargetConfig) WriteTargetConfigToWire(w *ConfigurationTarget) error {
  return nil
}


