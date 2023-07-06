package main

import (
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type gauge struct {
	name string
	val  float64
}
type counter struct {
	name string
	val  int
}

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second

func main() {
	var m runtime.MemStats

	go readRuntimeMetrics(&m)
	go pullMetrics(&m)
	for {

	}

}

func readRuntimeMetrics(m *runtime.MemStats) {
	for {
		runtime.ReadMemStats(m)
		time.Sleep(pollInterval)
	}

}
func pullMetrics(m *runtime.MemStats) {
	client := http.Client{Timeout: 100 * time.Millisecond}
	var RandomValue float64
	for {
		client.Post("http://127.0.0.1:8080/update/gauge/Alloc/"+strconv.FormatUint(m.Alloc, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/BuckHashSys/"+strconv.FormatUint(m.BuckHashSys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/Frees/"+strconv.FormatUint(m.Frees, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/GCCPUFraction/"+strconv.FormatFloat(m.GCCPUFraction, 'f', -1, 64), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/GCSys/"+strconv.FormatUint(m.GCSys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/HeapAlloc/"+strconv.FormatUint(m.HeapAlloc, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/HeapIdle/"+strconv.FormatUint(m.HeapIdle, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/HeapInuse/"+strconv.FormatUint(m.HeapInuse, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/HeapObjects/"+strconv.FormatUint(m.HeapObjects, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/HeapReleased/"+strconv.FormatUint(m.HeapReleased, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/HeapSys/"+strconv.FormatUint(m.HeapSys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/LastGC/"+strconv.FormatUint(m.LastGC, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/Lookups/"+strconv.FormatUint(m.Lookups, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/Mallocs/"+strconv.FormatUint(m.Mallocs, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/MSpanSys/"+strconv.FormatUint(m.MSpanSys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/MSpanInuse/"+strconv.FormatUint(m.MSpanInuse, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/MCacheSys/"+strconv.FormatUint(m.MCacheSys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/NextGC/"+strconv.FormatUint(m.NextGC, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/MCacheInuse/"+strconv.FormatUint(uint64(m.NumForcedGC), 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/NumGC/"+strconv.FormatUint(uint64(m.NumGC), 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/OtherSys/"+strconv.FormatUint(m.OtherSys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/PauseTotalNs/"+strconv.FormatUint(m.PauseTotalNs, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/StackInuse/"+strconv.FormatUint(m.StackInuse, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/StackSys/"+strconv.FormatUint(m.StackSys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/Sys/"+strconv.FormatUint(m.Sys, 10), "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/TotalAlloc/"+strconv.FormatUint(m.TotalAlloc, 10), "text/plain", nil)

		client.Post("http://127.0.0.1:8080/update/counter/PollCount/1", "text/plain", nil)
		client.Post("http://127.0.0.1:8080/update/gauge/RandomValue/"+strconv.FormatFloat(RandomValue, 'f', -1, 64), "text/plain", nil)
		time.Sleep(reportInterval)
	}

}

func sendGauge(g gauge, client http.Client) (*http.Response, error) {
	resp, err := client.Post("http://127.0.0.1:8080/update/gauge/"+g.name+strconv.FormatFloat(g.val, 'f', -1, 64), "text/plain", nil)
	return resp, err
}
