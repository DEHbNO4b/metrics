package agent

import (
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var pollInterval = 2 * time.Second
var reportInterval = 10 * time.Second

func ReadRuntimeMetrics(m *runtime.MemStats) {
	for {
		runtime.ReadMemStats(m)
		time.Sleep(pollInterval)
	}

}
func PullMetrics(m *runtime.MemStats) {
	url := "http://127.0.0.1:8080/update/"
	client := http.Client{Timeout: 1000 * time.Millisecond}
	// client := http.Client{}
	var RandomValue float64
	for {
		RandomValue = rand.Float64()
		go sendGauge(url+"Alloc/"+strconv.FormatUint(m.Alloc, 10), client)
		go sendGauge(url+"BuckHashSys/"+strconv.FormatUint(m.BuckHashSys, 10), client)
		go sendGauge(url+"Frees/"+strconv.FormatUint(m.Frees, 10), client)
		go sendGauge(url+"GCCPUFraction/"+strconv.FormatFloat(m.GCCPUFraction, 'f', -1, 64), client)
		go sendGauge(url+"GCSys/"+strconv.FormatUint(m.GCSys, 10), client)
		go sendGauge(url+"HeapAlloc/"+strconv.FormatUint(m.HeapAlloc, 10), client)
		go sendGauge(url+"HeapIdle/"+strconv.FormatUint(m.HeapIdle, 10), client)
		go sendGauge(url+"HeapInuse/"+strconv.FormatUint(m.HeapInuse, 10), client)
		go sendGauge(url+"HeapObjects/"+strconv.FormatUint(m.HeapObjects, 10), client)
		go sendGauge(url+"HeapReleased/"+strconv.FormatUint(m.HeapReleased, 10), client)
		go sendGauge(url+"HeapSys/"+strconv.FormatUint(m.HeapSys, 10), client)
		go sendGauge(url+"LastGC/"+strconv.FormatUint(m.LastGC, 10), client)
		go sendGauge(url+"Lookups/"+strconv.FormatUint(m.Lookups, 10), client)
		go sendGauge(url+"Mallocs/"+strconv.FormatUint(m.Mallocs, 10), client)
		go sendGauge(url+"MSpanSys/"+strconv.FormatUint(m.MSpanSys, 10), client)
		go sendGauge(url+"MSpanInuse/"+strconv.FormatUint(m.MSpanInuse, 10), client)
		go sendGauge(url+"MCacheSys/"+strconv.FormatUint(m.MCacheSys, 10), client)
		go sendGauge(url+"NextGC/"+strconv.FormatUint(m.NextGC, 10), client)
		go sendGauge(url+"NumForcedGC/"+strconv.FormatUint(uint64(m.NumForcedGC), 10), client)
		go sendGauge(url+"NumGC/"+strconv.FormatUint(uint64(m.NumGC), 10), client)
		go sendGauge(url+"OtherSys/"+strconv.FormatUint(m.OtherSys, 10), client)
		go sendGauge(url+"PauseTotalNs/"+strconv.FormatUint(m.PauseTotalNs, 10), client)
		go sendGauge(url+"StackInuse/"+strconv.FormatUint(m.StackInuse, 10), client)
		go sendGauge(url+"StackSys/"+strconv.FormatUint(m.StackSys, 10), client)
		go sendGauge(url+"Sys/"+strconv.FormatUint(m.Sys, 10), client)
		go sendGauge(url+"TotalAlloc/"+strconv.FormatUint(m.TotalAlloc, 10), client)
		go sendGauge(url+"RandomValue/"+strconv.FormatFloat(RandomValue, 'f', -1, 64), client)

		resp, _ := client.Post("http://127.0.0.1:8080/update/counter/PollCount/1", "text/plain", nil)
		resp.Body.Close()

		time.Sleep(reportInterval)
	}

}

func sendGauge(uri string, client http.Client) error {
	resp, err := client.Post(uri, "text/plain", nil)
	resp.Body.Close()
	return err
}
