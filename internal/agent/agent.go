package agent

import (
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second

func ReadRuntimeMetrics(m *runtime.MemStats) {
	var lock sync.Mutex
	for {
		lock.Lock()
		runtime.ReadMemStats(m)
		lock.Unlock()
		time.Sleep(pollInterval)
	}

}
func PullMetrics(m *runtime.MemStats) {
	url := "http://localhost:8080/update/gauge/"
	client := http.Client{Timeout: 10000 * time.Millisecond}
	// client := http.Client{}
	var RandomValue float64
	for {

		sendGauge(url+"Alloc/"+strconv.FormatUint(m.Alloc, 10), client)
		sendGauge(url+"BuckHashSys/"+strconv.FormatUint(m.BuckHashSys, 10), client)
		sendGauge(url+"Frees/"+strconv.FormatUint(m.Frees, 10), client)
		sendGauge(url+"GCCPUFraction/"+strconv.FormatFloat(m.GCCPUFraction, 'f', -1, 64), client)
		sendGauge(url+"GCSys/"+strconv.FormatUint(m.GCSys, 10), client)
		sendGauge(url+"HeapAlloc/"+strconv.FormatUint(m.HeapAlloc, 10), client)
		sendGauge(url+"HeapIdle/"+strconv.FormatUint(m.HeapIdle, 10), client)
		sendGauge(url+"HeapInuse/"+strconv.FormatUint(m.HeapInuse, 10), client)
		sendGauge(url+"HeapObjects/"+strconv.FormatUint(m.HeapObjects, 10), client)
		sendGauge(url+"HeapReleased/"+strconv.FormatUint(m.HeapReleased, 10), client)
		sendGauge(url+"HeapSys/"+strconv.FormatUint(m.HeapSys, 10), client)
		sendGauge(url+"LastGC/"+strconv.FormatUint(m.LastGC, 10), client)
		sendGauge(url+"Lookups/"+strconv.FormatUint(m.Lookups, 10), client)
		sendGauge(url+"Mallocs/"+strconv.FormatUint(m.Mallocs, 10), client)
		sendGauge(url+"MSpanSys/"+strconv.FormatUint(m.MSpanSys, 10), client)
		sendGauge(url+"MSpanInuse/"+strconv.FormatUint(m.MSpanInuse, 10), client)
		sendGauge(url+"MCacheSys/"+strconv.FormatUint(m.MCacheSys, 10), client)
		sendGauge(url+"NextGC/"+strconv.FormatUint(m.NextGC, 10), client)
		sendGauge(url+"NumForcedGC/"+strconv.FormatUint(uint64(m.NumForcedGC), 10), client)
		sendGauge(url+"NumGC/"+strconv.FormatUint(uint64(m.NumGC), 10), client)
		sendGauge(url+"OtherSys/"+strconv.FormatUint(m.OtherSys, 10), client)
		sendGauge(url+"PauseTotalNs/"+strconv.FormatUint(m.PauseTotalNs, 10), client)
		sendGauge(url+"StackInuse/"+strconv.FormatUint(m.StackInuse, 10), client)
		sendGauge(url+"StackSys/"+strconv.FormatUint(m.StackSys, 10), client)
		sendGauge(url+"Sys/"+strconv.FormatUint(m.Sys, 10), client)
		sendGauge(url+"TotalAlloc/"+strconv.FormatUint(m.TotalAlloc, 10), client)
		RandomValue = rand.Float64()
		sendGauge(url+"RandomValue/"+strconv.FormatFloat(RandomValue, 'f', -1, 64), client)

		resp, _ := client.Post("http://127.0.0.1:8080/update/counter/PollCount/1", "text/plain", nil)
		resp.Body.Close()

		time.Sleep(reportInterval)
	}

}

func sendGauge(uri string, client http.Client) error {
	resp, err := client.Post(uri, "text/plain", nil)
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return err
}
