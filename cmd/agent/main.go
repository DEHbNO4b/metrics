package main

import (
	"runtime"

	"github.com/DEHbNO4b/metrics/internal/agent"
)

func main() {
	// parseFlag()
	// a := agent.NewAgent(endpoint)
	// go a.ReadRuntimeMetrics(pollInterval)
	// go a.PullMetrics(reportInterval)
	// select {}
	parseFlag()
	var m runtime.MemStats
	go agent.ReadRuntimeMetrics(&m, pollInterval)
	go agent.PullMetrics(&m, reportInterval, endpoint)
	select {}
}
