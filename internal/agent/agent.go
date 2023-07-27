package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

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
			el.ReadValue(a.m)
			// reader, ok := metricReaders[el.ID]
			// if !ok {
			// 	logger.Log.Error("dont have metric with that name", zap.String("name", el.ID))
			// }
			// val := reader(a.m)
			// el.Value = &val
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
