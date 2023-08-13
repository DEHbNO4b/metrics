package agent

import (
	"bytes"
	"compress/gzip"
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
		// for _, el := range a.gauges {
		// 	el.ReadValue(a.m)
		// 	go a.sendMetric(el)
		// }
		d := int64(1)
		// go a.sendMetric(data.Metrics{ID: "PollCount", MType: "counter", Delta: &d})
		a.gauges = append(a.gauges, data.Metrics{ID: "PollCount", MType: "counter", Delta: &d})
		go a.sendMetrics(a.gauges)
		time.Sleep(reportInterval)
	}

}
func (a Agent) sendMetrics(metrics []data.Metrics) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&metrics)
	if err != nil {
		logger.Log.Info("unable to encode metric", zap.String("err: ", err.Error()))
		return
	}
	compressed := bytes.Buffer{}
	compressor, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	compressor.Write(buf.Bytes())
	compressor.Close()
	req, err := http.NewRequest(http.MethodPost, a.url+"/update/", &compressed) // (1)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-encoding", "gzip")
	req.Header.Add("Accept-encoding", "gzip")
	resp, err := a.client.Do(req)
	if err != nil {
		logger.Log.Error("request returned err ", zap.Error(err))
		return
	}
	resp.Body.Close()
}
func (a Agent) sendMetric(m data.Metrics) {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(&m)
	if err != nil {
		logger.Log.Info("unable to encode metric", zap.String("err: ", err.Error()))
		return
	}
	compressed := bytes.Buffer{}
	compressor, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	// fmt.Println(buf.Bytes())
	compressor.Write(buf.Bytes())
	compressor.Close()

	req, err := http.NewRequest(http.MethodPost, a.url+"/update/", &compressed) // (1)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-encoding", "gzip")
	req.Header.Add("Accept-encoding", "gzip")
	resp, err := a.client.Do(req)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		return
	}
	resp.Body.Close()
}
