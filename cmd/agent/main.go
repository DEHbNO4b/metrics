// package agent provides finctions for collecting metrics from *runtime.MemStats.
package main

import (
	"github.com/DEHbNO4b/metrics/internal/agent"
)

func main() {
	parseFlag()
	a := agent.NewAgent(endpoint)
	go a.ReadRuntimeMetrics(pollInterval)
	go a.PullMetrics(reportInterval, key, crypto)
	select {}

}
