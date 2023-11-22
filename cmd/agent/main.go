// package agent provides finctions for collecting metrics from *runtime.MemStats.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DEHbNO4b/metrics/internal/agent"
	"github.com/DEHbNO4b/metrics/internal/config"
)

func main() {
	// parseFlag()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.GetAgentCfg()
	a := agent.NewAgent(cfg.Adress)
	go a.ReadRuntimeMetrics(ctx, cfg.PollInterval)
	go a.PullMetrics(ctx, cfg.ReportInterval, cfg.HashKey, cfg.CryptoKey)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigint

		cancel()
	}()

	<-ctx.Done()
	time.Sleep(5 * time.Second)
}
