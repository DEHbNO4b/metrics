package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	"github.com/go-chi/chi/v5"
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
	fmt.Println(req.URL.Path)
	url, _ := strings.CutPrefix(req.URL.Path, "/update/")
	urlValues := strings.Split(url, "/")
	if len(urlValues) < 3 {
		http.Error(w, "bad request", http.StatusNotFound)
		return
	}
	if urlValues[0] == "" {
		http.Error(w, "bad request", http.StatusNotFound)
		return
	}
	if urlValues[1] == "" {
		http.Error(w, "bad request", http.StatusNotFound)
		return
	}
	if urlValues[0] != "counter" && urlValues[0] != "gauge" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// fmt.Println(url)

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
	val, err := strconv.ParseInt(urlValues[1], 10, 64)
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	// ms.counter[urlValues[0]] += number
	ms.MemStorage.SetCounter(data.Counter{Name: urlValues[0], Val: val})
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
