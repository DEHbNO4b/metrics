package agent

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

// var gauges = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapReleased", "HeapObjects", "HeapSys", "LastGC",
// 	"Lookups", "MCacheInuse", "Mallocs", "MSpanSys", "MSpanInuse", "MCacheSys", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}
// var counters = []string{"PollCount"}
// var ga map[string]func() float32 = make(map[string]func() float32)

type Agent struct {
	m      runtime.MemStats
	client http.Client
	url    string
}

func NewAgent(endpoint string) Agent {
	var me runtime.MemStats
	cl := http.Client{Timeout: 1000 * time.Millisecond}
	u := "http://" + endpoint
	a := Agent{m: me, client: cl, url: u}
	return a
}
func (a Agent) ReadRuntimeMetrics(interval int) {
	var pollInterval = time.Duration(interval) * time.Second
	for {
		runtime.ReadMemStats(&a.m)
		time.Sleep(pollInterval)
	}

}
func (a Agent) PullMetrics(interval int) {
	var reportInterval = time.Duration(interval) * time.Second

	var RandomValue float64
	for {
		// for _, el := range gauges {
		// 	a.sendMetric(el)
		// }
		al := float64(a.m.Alloc)
		go a.sendMetric(data.Metrics{ID: "Alloc", MType: "gauge", Value: &al})
		bh := float64(a.m.BuckHashSys)
		go a.sendMetric(data.Metrics{ID: "BuckHashSys", MType: "gauge", Value: &bh})
		f := float64(a.m.Frees)
		go a.sendMetric(data.Metrics{ID: "Frees", MType: "gauge", Value: &f})
		gf := float64(a.m.GCCPUFraction)
		go a.sendMetric(data.Metrics{ID: "GCCPUFraction", MType: "gauge", Value: &gf})
		gs := float64(a.m.GCSys)
		go a.sendMetric(data.Metrics{ID: "GCSys", MType: "gauge", Value: &gs})
		ha := float64(a.m.HeapAlloc)
		go a.sendMetric(data.Metrics{ID: "HeapAlloc", MType: "gauge", Value: &ha})
		hi := float64(a.m.HeapIdle)
		go a.sendMetric(data.Metrics{ID: "HeapIdle", MType: "gauge", Value: &hi})
		hin := float64(a.m.HeapInuse)
		go a.sendMetric(data.Metrics{ID: "HeapInuse", MType: "gauge", Value: &hin})
		ho := float64(a.m.HeapObjects)
		go a.sendMetric(data.Metrics{ID: "HeapObjects", MType: "gauge", Value: &ho})
		hr := float64(a.m.HeapReleased)
		go a.sendMetric(data.Metrics{ID: "HeapReleased", MType: "gauge", Value: &hr})
		hs := float64(a.m.HeapSys)
		go a.sendMetric(data.Metrics{ID: "HeapSys", MType: "gauge", Value: &hs})
		lgc := float64(a.m.LastGC)
		go a.sendMetric(data.Metrics{ID: "LastGC", MType: "gauge", Value: &lgc})
		l := float64(a.m.Lookups)
		go a.sendMetric(data.Metrics{ID: "Lookups", MType: "gauge", Value: &l})
		mi := float64(a.m.MCacheInuse)
		go a.sendMetric(data.Metrics{ID: "MCacheInuse", MType: "gauge", Value: &mi})
		m := float64(a.m.Mallocs)
		go a.sendMetric(data.Metrics{ID: "Alloc", MType: "gauge", Value: &m})
		ms := float64(a.m.MSpanSys)
		go a.sendMetric(data.Metrics{ID: "MSpanSys", MType: "gauge", Value: &ms})
		msi := float64(a.m.MSpanInuse)
		go a.sendMetric(data.Metrics{ID: "MSpanInuse", MType: "gauge", Value: &msi})
		mcs := float64(a.m.MCacheSys)
		go a.sendMetric(data.Metrics{ID: "MCacheSys", MType: "gauge", Value: &mcs})
		ngc := float64(a.m.NextGC)
		go a.sendMetric(data.Metrics{ID: "NextGC", MType: "gauge", Value: &ngc})
		nfg := float64(a.m.NumForcedGC)
		go a.sendMetric(data.Metrics{ID: "NumForcedGC", MType: "gauge", Value: &nfg})
		ng := float64(a.m.NumGC)
		go a.sendMetric(data.Metrics{ID: "NumGC", MType: "gauge", Value: &ng})
		os := float64(a.m.OtherSys)
		go a.sendMetric(data.Metrics{ID: "OtherSys", MType: "gauge", Value: &os})
		pt := float64(a.m.PauseTotalNs)
		go a.sendMetric(data.Metrics{ID: "PauseTotalNs", MType: "gauge", Value: &pt})
		si := float64(a.m.StackInuse)
		go a.sendMetric(data.Metrics{ID: "StackInuse", MType: "gauge", Value: &si})
		ss := float64(a.m.StackSys)
		go a.sendMetric(data.Metrics{ID: "StackSys", MType: "gauge", Value: &ss})
		s := float64(a.m.Sys)
		go a.sendMetric(data.Metrics{ID: "Sys", MType: "gauge", Value: &s})
		ta := float64(a.m.TotalAlloc)
		go a.sendMetric(data.Metrics{ID: "TotalAlloc", MType: "gauge", Value: &ta})

		RandomValue = rand.Float64()
		go a.sendMetric(data.Metrics{ID: "RandomValue", MType: "gauge", Value: &RandomValue})

		d := int64(1)
		go a.sendMetric(data.Metrics{ID: "PollCount", MType: "counter", Delta: &d})

		time.Sleep(reportInterval)
	}

}

func (a Agent) sendMetric(m data.Metrics) error {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&m)
	if err != nil {
		logger.Log.Info("unable ti encode metric", zap.String("err: ", err.Error()))
		return err
	}
	resp, err := a.client.Post(a.url+"/update/", "application/json", &buf)

	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
