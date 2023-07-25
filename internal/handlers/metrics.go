package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

type Metrics struct {
	MemStorage interfaces.MetricsStorage
}

func NewMetrics(m interfaces.MetricsStorage) Metrics {
	ms := Metrics{MemStorage: m}
	return ms
}

func (ms *Metrics) SetMetrics(w http.ResponseWriter, req *http.Request) {
	// fmt.Println("in set metrics")
	m := data.Metrics{}
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&m)
	if err != nil {
		logger.Log.Info("unable to decode json", zap.String("err", err.Error()))
	}
	switch m.MType {
	case "gauge":
		ms.MemStorage.SetGauge(data.Gauge{Name: m.ID, Val: *m.Value})
	case "counter":
		ms.MemStorage.SetCounter(data.Counter{Name: m.ID, Val: *m.Delta})
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (ms *Metrics) SetGauge(w http.ResponseWriter, req *http.Request) {
	url, _ := strings.CutPrefix(req.URL.Path, "/update/gauge/")
	urlValues := strings.Split(url, "/")

	val, err := strconv.ParseFloat(urlValues[1], 64)
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	ms.MemStorage.SetGauge(data.Gauge{Name: urlValues[0], Val: val})
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (ms *Metrics) SetCounter(w http.ResponseWriter, req *http.Request) {

	url, _ := strings.CutPrefix(req.URL.Path, "/update/counter/")
	urlValues := strings.Split(url, "/")
	val, err := strconv.ParseInt(urlValues[1], 10, 64)
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	ms.MemStorage.SetCounter(data.Counter{Name: urlValues[0], Val: val})
	// w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))

}
func (ms *Metrics) GetMetrics(w http.ResponseWriter, r *http.Request) {

	//m.MemStorage.SetGauge(data.Gauge{Name: "qwe", Val: 234234})
	const formbegin = `<html>
			<head>
				<title></title>
			</head>
				<body>
			`
	const formend = `
				</body>
			</html>`
	metrics := ms.MemStorage.GetMetrics()

	io.WriteString(w, formbegin)
	io.WriteString(w, strings.Join(metrics, ", "))
	io.WriteString(w, formend)
}
func (ms *Metrics) GetMetric(w http.ResponseWriter, r *http.Request) {
	m := data.Metrics{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&m)
	switch m.MType {
	case "gauge":
		g, err := ms.MemStorage.GetGauge(m.ID)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
		}
		m = data.Metrics{ID: g.Name, MType: "gauge", Value: &g.Val}

	case "counter":
		c, err := ms.MemStorage.GetCounter(m.ID)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
		}
		m = data.Metrics{ID: c.Name, MType: "counter", Delta: &c.Val}
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(&m)
}
