package main

import (
	"github.com/DEHbNO4b/metrics/internal/agent"
)

func main() {
	parseFlag()
	a := agent.NewAgent(endpoint)
	go a.ReadRuntimeMetrics(pollInterval)
	go a.PullMetrics(reportInterval)
	select {}
	// parseFlag()
	// // gauges := data.NewGauges()
	// var m runtime.MemStats
	// go agent.ReadRuntimeMetrics(&m, pollInterval)
	// go agent.PullMetrics( &m, reportInterval, endpoint)
	// select {}
}
