package agent

import (
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func ReadRuntimeMetrics(m *runtime.MemStats, interval int) {
	var pollInterval = time.Duration(interval) * time.Second

	var lock sync.Mutex
	for {
		lock.Lock()
		runtime.ReadMemStats(m)
		lock.Unlock()
		time.Sleep(pollInterval)
	}

}

func PullMetrics(m *runtime.MemStats, interval int, endpoint string) {
	var reportInterval = time.Duration(interval) * time.Second

	url := "http://" + endpoint
	urlUpdateGauge := url + "/update/gauge/"
	urlUpdateCounter := url + "/update/counter/"

	client := http.Client{Timeout: 1000 * time.Millisecond}
	var RandomValue float64
	for {

		go sendMetric(urlUpdateGauge+"Alloc/"+strconv.FormatUint(m.Alloc, 10), client)
		go sendMetric(urlUpdateGauge+"BuckHashSys/"+strconv.FormatUint(m.BuckHashSys, 10), client)
		go sendMetric(urlUpdateGauge+"Frees/"+strconv.FormatUint(m.Frees, 10), client)
		go sendMetric(urlUpdateGauge+"GCCPUFraction/"+strconv.FormatFloat(m.GCCPUFraction, 'f', -1, 64), client)
		go sendMetric(urlUpdateGauge+"GCSys/"+strconv.FormatUint(m.GCSys, 10), client)
		go sendMetric(urlUpdateGauge+"HeapAlloc/"+strconv.FormatUint(m.HeapAlloc, 10), client)
		go sendMetric(urlUpdateGauge+"HeapIdle/"+strconv.FormatUint(m.HeapIdle, 10), client)
		go sendMetric(urlUpdateGauge+"HeapInuse/"+strconv.FormatUint(m.HeapInuse, 10), client)
		go sendMetric(urlUpdateGauge+"HeapObjects/"+strconv.FormatUint(m.HeapObjects, 10), client)
		go sendMetric(urlUpdateGauge+"HeapReleased/"+strconv.FormatUint(m.HeapReleased, 10), client)
		go sendMetric(urlUpdateGauge+"HeapSys/"+strconv.FormatUint(m.HeapSys, 10), client)
		go sendMetric(urlUpdateGauge+"LastGC/"+strconv.FormatUint(m.LastGC, 10), client)
		go sendMetric(urlUpdateGauge+"Lookups/"+strconv.FormatUint(m.Lookups, 10), client)
		go sendMetric(urlUpdateGauge+"MCacheInuse/"+strconv.FormatUint(m.MCacheInuse, 10), client)
		go sendMetric(urlUpdateGauge+"Mallocs/"+strconv.FormatUint(m.Mallocs, 10), client)
		go sendMetric(urlUpdateGauge+"MSpanSys/"+strconv.FormatUint(m.MSpanSys, 10), client)
		go sendMetric(urlUpdateGauge+"MSpanInuse/"+strconv.FormatUint(m.MSpanInuse, 10), client)
		go sendMetric(urlUpdateGauge+"MCacheSys/"+strconv.FormatUint(m.MCacheSys, 10), client)
		go sendMetric(urlUpdateGauge+"NextGC/"+strconv.FormatUint(m.NextGC, 10), client)
		go sendMetric(urlUpdateGauge+"NumForcedGC/"+strconv.FormatUint(uint64(m.NumForcedGC), 10), client)
		go sendMetric(urlUpdateGauge+"NumGC/"+strconv.FormatUint(uint64(m.NumGC), 10), client)
		go sendMetric(urlUpdateGauge+"OtherSys/"+strconv.FormatUint(m.OtherSys, 10), client)
		go sendMetric(urlUpdateGauge+"PauseTotalNs/"+strconv.FormatUint(m.PauseTotalNs, 10), client)
		go sendMetric(urlUpdateGauge+"StackInuse/"+strconv.FormatUint(m.StackInuse, 10), client)
		go sendMetric(urlUpdateGauge+"StackSys/"+strconv.FormatUint(m.StackSys, 10), client)
		go sendMetric(urlUpdateGauge+"Sys/"+strconv.FormatUint(m.Sys, 10), client)
		go sendMetric(urlUpdateGauge+"TotalAlloc/"+strconv.FormatUint(m.TotalAlloc, 10), client)

		RandomValue = rand.Float64()
		go sendMetric(urlUpdateGauge+"RandomValue/"+strconv.FormatFloat(RandomValue, 'f', -1, 64), client)
		go sendMetric(urlUpdateCounter+"PollCount/1", client)

		time.Sleep(reportInterval)
	}

}

func sendMetric(uri string, client http.Client) error {
	resp, err := client.PostForm(uri, nil)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}
