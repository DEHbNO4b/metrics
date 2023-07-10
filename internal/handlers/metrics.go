package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
)

type Metrics struct {
	MemStorage interfaces.MetricsStorage
	// gauge   map[string]float64
	// counter map[string]int
}

func NewMetrics(m interfaces.MetricsStorage) Metrics {
	// g := make(map[string]float64)
	// c := make(map[string]int)
	ms := Metrics{MemStorage: m}
	return ms
}

func (ms *Metrics) SetMetrics(w http.ResponseWriter, req *http.Request) {
	url, _ := strings.CutPrefix(req.URL.Path, "/update/")
	// fmt.Println(url)
	urlValues := strings.Split(url, "/")

	switch urlValues[0] {
	case "gauge":
		ms.SetGauge(w, req)
	case "counter":
		ms.SetCounter(w, req)
	default:
		{
			http.Error(w, "wrong metric value", http.StatusBadRequest)
			return
		}
	}
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (ms *Metrics) SetCounter(w http.ResponseWriter, req *http.Request) {

	url, _ := strings.CutPrefix(req.URL.Path, "/update/counter/")
	urlValues := strings.Split(url, "/")
	val, err := strconv.Atoi(urlValues[1])
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	// ms.counter[urlValues[0]] += number
	ms.MemStorage.SetCounter(data.Counter{Name: urlValues[0], Val: val})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))

}
