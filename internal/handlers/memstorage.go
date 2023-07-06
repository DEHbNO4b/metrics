package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int
}

func NewMemStorage() *MemStorage {
	g := make(map[string]float64)
	c := make(map[string]int)
	ms := MemStorage{gauge: g, counter: c}
	return &ms
}

func (ms *MemStorage) SetMetrics(w http.ResponseWriter, req *http.Request) {
	url, _ := strings.CutPrefix(req.URL.Path, "/update/")
	fmt.Println(url)
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

func (ms *MemStorage) SetGauge(w http.ResponseWriter, req *http.Request) {
	url, _ := strings.CutPrefix(req.URL.Path, "/update/gauge/")
	urlValues := strings.Split(url, "/")

	number, err := strconv.ParseFloat(urlValues[1], 64)
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	ms.gauge[urlValues[0]] = number
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("set gauge"))
}

func (ms *MemStorage) SetCounter(w http.ResponseWriter, req *http.Request) {

	url, _ := strings.CutPrefix(req.URL.Path, "/update/counter/")
	urlValues := strings.Split(url, "/")
	number, err := strconv.Atoi(urlValues[1])
	if err != nil {
		http.Error(w, "wrong metric value", http.StatusBadRequest)
		return
	}
	ms.counter[urlValues[0]] += number
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Set counter"))

}
