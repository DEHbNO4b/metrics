package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

var metricReaders = map[string]func(m *runtime.MemStats) float64{
	"Alloc": func(m *runtime.MemStats) float64 {
		fmt.Println(m.Alloc)
		return float64(m.Alloc)
	},
	"BuckHashSys": func(m *runtime.MemStats) float64 {
		fmt.Println(m.BuckHashSys)
		return float64(m.BuckHashSys)
	},
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

type Agent struct {
	m      *runtime.MemStats
	client http.Client
	url    string
	gauges []data.Metrics
}

func NewAgent(endpoint string) Agent {
	var m runtime.MemStats
	cl := http.Client{Timeout: 1000 * time.Millisecond}
	u := "http://" + endpoint
	a := Agent{m: &m, client: cl, url: u, gauges: data.NewGauges()}
	return a
}

func (a Agent) ReadRuntimeMetrics(interval int) {
	var pollInterval = time.Duration(interval) * time.Second
	for {
		runtime.ReadMemStats(a.m)
		time.Sleep(pollInterval)
	}

}
func (a Agent) PullMetrics(interval int) {
	var reportInterval = time.Duration(interval) * time.Second
	for {
		for _, el := range a.gauges {
			reader, ok := metricReaders[el.ID]
			if !ok {
				logger.Log.Error("dont have metric with that name", zap.String("name", el.ID))
			}
			val := reader(a.m)
			el.Value = &val
			go a.sendMetric(el)
		}
		d := int64(1)
		go a.sendMetric(data.Metrics{ID: "PollCount", MType: "counter", Delta: &d})

		time.Sleep(reportInterval)
	}

}

func (a Agent) sendMetric(m data.Metrics) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&m)
	fmt.Println(buf.String())
	if err != nil {
		logger.Log.Info("unable to encode metric", zap.String("err: ", err.Error()))
		return
	}
	resp, err := a.client.Post(a.url+"/update/", "application/json", &buf)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	resp.Body.Close()
}
