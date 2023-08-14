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
	expert interfaces.MetricsStorage
	Pinger interfaces.Pinger
}

func NewMetrics(m interfaces.MetricsStorage) Metrics {
	ms := Metrics{expert: m}
	return ms
}

func (ms *Metrics) SetMetricJSON(w http.ResponseWriter, req *http.Request) {
	m := data.Metrics{}
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&m)
	if err != nil {
		logger.Log.Info("unable to decode json from req.Body", zap.String("err", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = ms.expert.SetMetric(m)
	if err != nil {
		http.Error(w, "unable to set metrics to RAM", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (ms *Metrics) SetMetricsJSON(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("in set metrics handler")
	metrics := make([]data.Metrics, 0, 30)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&metrics)
	if err != nil {
		logger.Log.Info("unable to decode json from req.Body", zap.String("err", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	for _, metric := range metrics {
		err = ms.expert.SetMetric(metric)
		if err != nil {
			http.Error(w, "unable to set metrics to RAM", http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (ms *Metrics) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	m := data.Metrics{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&m)
	if err != nil {
		http.Error(w, "unable to decode teq body", http.StatusBadRequest)
		return
	}

	m, err = ms.expert.GetMetric(m)
	if err != nil && err == interfaces.ErrWrongType {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err != nil && err == interfaces.ErrNotContains {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(&m)
}

func (ms *Metrics) GetMetrics(w http.ResponseWriter, r *http.Request) {
	const formbegin = `<html><head><title></title></head><body>`
	const formend = `</body></html>`
	metrics := ms.expert.GetMetrics()
	m := make([]string, 0, 40)
	for _, val := range metrics {
		switch val.MType {
		case "gauge":
			m = append(m, val.ID+":"+strconv.FormatFloat(*val.Value, 'f', -1, 64))
		case "counter":
			m = append(m, val.ID+":"+strconv.FormatInt(*val.Delta, 10))
		}

	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, formbegin)
	io.WriteString(w, strings.Join(m, ", "))
	io.WriteString(w, formend)
}

func (ms *Metrics) SetMetricsURL(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("in set metrics url")
	url, _ := strings.CutPrefix(r.URL.Path, "/update/")
	urlValues := strings.Split(url, "/")
	if urlValues[1] == "" {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	m := data.Metrics{}
	m.MType = urlValues[0]
	m.ID = urlValues[1]
	switch urlValues[0] {
	case "gauge":
		val, err := strconv.ParseFloat(urlValues[2], 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		m.Value = &val
	case "counter":
		del, err := strconv.ParseInt(urlValues[2], 10, 64)
		if err != nil {
			http.Error(w, "Wrong metric value", http.StatusBadRequest)
			return
		}
		m.Delta = &del
	default:
		{
			http.Error(w, "Wrong metric type", http.StatusBadRequest)
			return
		}
	}
	err := ms.expert.SetMetric(m)
	if err != nil {
		logger.Log.Sugar().Error(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ms *Metrics) GetMetricURL(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	metric := data.Metrics{ID: name, MType: t}
	m, err := ms.expert.GetMetric(metric)
	if err != nil && err == interfaces.ErrWrongType {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err != nil && err == interfaces.ErrNotContains {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	data := ""
	switch m.MType {
	case "gauge":
		data = strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case "counter":
		data = strconv.FormatInt(*m.Delta, 10)
	default:
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.Write([]byte(data))
}
