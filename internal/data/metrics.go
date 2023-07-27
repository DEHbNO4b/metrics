package data

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

var gauges = []Metrics{
	{
		ID:    "Alloc",
		MType: "gauge",
	},
	{
		ID:    "BuckHashSys",
		MType: "gauge",
	},
	{
		ID:    "Frees",
		MType: "gauge",
	},
	{
		ID:    "GCCPUFraction",
		MType: "gauge",
	},
	{
		ID:    "GCSys",
		MType: "gauge",
	},
	{
		ID:    "HeapAlloc",
		MType: "gauge",
	},
	{
		ID:    "HeapIdle",
		MType: "gauge",
	},
	{
		ID:    "HeapInuse",
		MType: "gauge",
	},
	{
		ID:    "HeapObjects",
		MType: "gauge",
	},
	{
		ID:    "HeapReleased",
		MType: "gauge",
	},
	{
		ID:    "HeapSys",
		MType: "gauge",
	},
	{
		ID:    "LastGC",
		MType: "gauge",
	},
	{
		ID:    "Lookups",
		MType: "gauge",
	},
	{
		ID:    "MCacheInuse",
		MType: "gauge",
	},
	{
		ID:    "MCacheSys",
		MType: "gauge",
	},
	{
		ID:    "MSpanInuse",
		MType: "gauge",
	},
	{
		ID:    "MSpanSys",
		MType: "gauge",
	},
	{
		ID:    "Mallocs",
		MType: "gauge",
	},
	{
		ID:    "NextGC",
		MType: "gauge",
	},
	{
		ID:    "NumForcedGC",
		MType: "gauge",
	},
	{
		ID:    "NumGC",
		MType: "gauge",
	},
	{
		ID:    "OtherSys",
		MType: "gauge",
	},
	{
		ID:    "PauseTotalNs",
		MType: "gauge",
	},
	{
		ID:    "StackInuse",
		MType: "gauge",
	},
	{
		ID:    "StackSys",
		MType: "gauge",
	},
	{
		ID:    "Sys",
		MType: "gauge",
	},
	{
		ID:    "TotalAlloc",
		MType: "gauge",
	},
}

func NewGauges() []Metrics {
	for _, el := range gauges {
		var val float64 = 0
		el.Value = &val
	}
	return gauges
}
