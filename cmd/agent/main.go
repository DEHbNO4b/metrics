package main

import (
	"runtime"

	"github.com/DEHbNO4b/metrics/internal/agent"
)

func main() {
	parseFlag()
	var m runtime.MemStats
	go agent.ReadRuntimeMetrics(&m, pollInterval)
	go agent.PullMetrics(&m, reportInterval, endpoint)
	select {}
}
