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
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Metrics struct {
	MemStorage interfaces.MetricsStorage
	Pinger     interfaces.Pinger
}

func NewMetrics(m interfaces.MetricsStorage) Metrics {
	ms := Metrics{MemStorage: m}
	return ms
}

func (ms *Metrics) SetMetricsJSON(w http.ResponseWriter, req *http.Request) {
	m := data.Metrics{}
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&m)
	if err != nil {
		logger.Log.Info("unable to decode json", zap.String("err", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = ms.MemStorage.SetMetric(m)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (ms *Metrics) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
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

func (ms *Metrics) GetMetrics(w http.ResponseWriter, r *http.Request) {
	const formbegin = `<html><head><title></title></head><body>`
	const formend = `</body></html>`
	metrics := ms.MemStorage.GetMetrics()
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, formbegin)
	io.WriteString(w, strings.Join(metrics, ", "))
	io.WriteString(w, formend)
}

func (ms *Metrics) SetMetricsURL(w http.ResponseWriter, req *http.Request) {
	url, _ := strings.CutPrefix(req.URL.Path, "/update/")
	urlValues := strings.Split(url, "/")
	if len(urlValues) < 3 {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if urlValues[0] == "" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if urlValues[1] == "" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if urlValues[0] != "counter" && urlValues[0] != "gauge" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	switch urlValues[0] {
	case "gauge":
		ms.SetGaugeURL(w, req)
	case "counter":
		ms.SetCounterURL(w, req)
	default:
		{
			http.Error(w, "Wrong metric type", http.StatusBadRequest)
			return
		}
	}
}
func (ms *Metrics) SetGaugeURL(w http.ResponseWriter, req *http.Request) {
	url, _ := strings.CutPrefix(req.URL.Path, "/update/gauge/")
	urlValues := strings.Split(url, "/")

	val, err := strconv.ParseFloat(urlValues[1], 64)
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	ms.MemStorage.SetGauge(data.Gauge{Name: urlValues[0], Val: val})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (ms *Metrics) SetCounterURL(w http.ResponseWriter, req *http.Request) {

	url, _ := strings.CutPrefix(req.URL.Path, "/update/counter/")
	urlValues := strings.Split(url, "/")
	val, err := strconv.ParseInt(urlValues[1], 10, 64)
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	ms.MemStorage.SetCounter(data.Counter{Name: urlValues[0], Val: val})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))

}

func (ms *Metrics) GetMetricURL(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	data := ""
	switch t {
	case "gauge":
		g, err := ms.MemStorage.GetGauge(name)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
		}
		data = strconv.FormatFloat(g.Val, 'f', -1, 64)
	case "counter":
		c, err := ms.MemStorage.GetCounter(name)
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
		}
		data = strconv.FormatInt(c.Val, 10)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
	w.Write([]byte(data))
}
func (ms *Metrics) PingDb(w http.ResponseWriter, r *http.Request) {
	err := ms.Pinger.Ping()
	if err != nil {
		http.Error(w, "db disconeccted", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
