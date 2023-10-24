// package data contains definition of main models.
package data

import (
	"math/rand"
	"runtime"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

// Metrics struct is a metric model.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewMetric return a new Metric.
func NewMetric() Metrics {
	var delta int64 = 0
	var value float64 = 0
	return Metrics{Delta: &delta, Value: &value}
}

// ReadValue reads value of specific runtime metric from *runtime.MemStats.
func (m *Metrics) ReadValue(mem *runtime.MemStats) {
	reader, ok := metricReaders[m.ID]
	if !ok {
		logger.Log.Error("dont have metric with that name", zap.String("name", m.ID))
	}
	val := reader(mem)
	m.Value = &val
}

var metricReaders = map[string]func(m *runtime.MemStats) float64{
	"Alloc":         func(m *runtime.MemStats) float64 { return float64(m.Alloc) },
	"BuckHashSys":   func(m *runtime.MemStats) float64 { return float64(m.BuckHashSys) },
	"Frees":         func(m *runtime.MemStats) float64 { return float64(m.Frees) },
	"GCCPUFraction": func(m *runtime.MemStats) float64 { return float64(m.GCCPUFraction) },
	"GCSys":         func(m *runtime.MemStats) float64 { return float64(m.GCSys) },
	"HeapAlloc":     func(m *runtime.MemStats) float64 { return float64(m.HeapAlloc) },
	"HeapIdle":      func(m *runtime.MemStats) float64 { return float64(m.HeapIdle) },
	"HeapInuse":     func(m *runtime.MemStats) float64 { return float64(m.HeapInuse) },
	"HeapObjects":   func(m *runtime.MemStats) float64 { return float64(m.HeapObjects) },
	"HeapReleased":  func(m *runtime.MemStats) float64 { return float64(m.HeapReleased) },
	"HeapSys":       func(m *runtime.MemStats) float64 { return float64(m.HeapSys) },
	"LastGC":        func(m *runtime.MemStats) float64 { return float64(m.LastGC) },
	"Lookups":       func(m *runtime.MemStats) float64 { return float64(m.Lookups) },
	"MCacheInuse":   func(m *runtime.MemStats) float64 { return float64(m.MCacheInuse) },
	"MCacheSys":     func(m *runtime.MemStats) float64 { return float64(m.MCacheSys) },
	"MSpanInuse":    func(m *runtime.MemStats) float64 { return float64(m.MSpanInuse) },
	"MSpanSys":      func(m *runtime.MemStats) float64 { return float64(m.MSpanSys) },
	"Mallocs":       func(m *runtime.MemStats) float64 { return float64(m.Mallocs) },
	"NextGC":        func(m *runtime.MemStats) float64 { return float64(m.NextGC) },
	"NumForcedGC":   func(m *runtime.MemStats) float64 { return float64(m.NumForcedGC) },
	"NumGC":         func(m *runtime.MemStats) float64 { return float64(m.NumGC) },
	"OtherSys":      func(m *runtime.MemStats) float64 { return float64(m.OtherSys) },
	"PauseTotalNs":  func(m *runtime.MemStats) float64 { return float64(m.PauseTotalNs) },
	"StackInuse":    func(m *runtime.MemStats) float64 { return float64(m.StackInuse) },
	"StackSys":      func(m *runtime.MemStats) float64 { return float64(m.StackSys) },
	"Sys":           func(m *runtime.MemStats) float64 { return float64(m.Sys) },
	"TotalAlloc":    func(m *runtime.MemStats) float64 { return float64(m.TotalAlloc) },
	"RandomValue":   func(m *runtime.MemStats) float64 { return rand.Float64() },
}
var gauges = []Metrics{
	{ID: "Alloc", MType: "gauge"},
	{ID: "BuckHashSys", MType: "gauge"},
	{ID: "Frees", MType: "gauge"},
	{ID: "GCCPUFraction", MType: "gauge"},
	{ID: "GCSys", MType: "gauge"},
	{ID: "HeapAlloc", MType: "gauge"},
	{ID: "HeapIdle", MType: "gauge"},
	{ID: "HeapInuse", MType: "gauge"},
	{ID: "HeapObjects", MType: "gauge"},
	{ID: "HeapReleased", MType: "gauge"},
	{ID: "HeapSys", MType: "gauge"},
	{ID: "LastGC", MType: "gauge"},
	{ID: "Lookups", MType: "gauge"},
	{ID: "MCacheInuse", MType: "gauge"},
	{ID: "MCacheSys", MType: "gauge"},
	{ID: "MSpanInuse", MType: "gauge"},
	{ID: "MSpanSys", MType: "gauge"},
	{ID: "Mallocs", MType: "gauge"},
	{ID: "NextGC", MType: "gauge"},
	{ID: "NumForcedGC", MType: "gauge"},
	{ID: "NumGC", MType: "gauge"},
	{ID: "OtherSys", MType: "gauge"},
	{ID: "PauseTotalNs", MType: "gauge"},
	{ID: "StackInuse", MType: "gauge"},
	{ID: "StackSys", MType: "gauge"},
	{ID: "Sys", MType: "gauge"},
	{ID: "TotalAlloc", MType: "gauge"},
	{ID: "RandomValue", MType: "gauge"},
}

func NewGauges() []Metrics {
	return gauges
}
