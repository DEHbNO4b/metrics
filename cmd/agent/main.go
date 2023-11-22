// package agent provides finctions for collecting metrics from *runtime.MemStats.
package main

import (
	"github.com/DEHbNO4b/metrics/internal/agent"
	"github.com/DEHbNO4b/metrics/internal/config"
)

func main() {
	// parseFlag()
	cfg := config.GetAgentCfg()
	a := agent.NewAgent(cfg.Adress)
	go a.ReadRuntimeMetrics(cfg.PollInterval)
	go a.PullMetrics(cfg.ReportInterval, cfg.HashKey, cfg.CryptoKey)
	select {}

}
