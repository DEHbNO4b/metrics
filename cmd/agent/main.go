package main

import (
	"runtime"
	"sync"

	"github.com/DEHbNO4b/metrics/internal/agent"
)

func main() {
	parseFlag()
	var m runtime.MemStats

	go agent.ReadRuntimeMetrics(&m, pollInterval)
	go agent.PullMetrics(&m, reportInterval, endpoint)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
