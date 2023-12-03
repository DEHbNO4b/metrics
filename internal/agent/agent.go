// packag agent provides automatic collection of metrics and sending them to the server.
package agent

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/DEHbNO4b/metrics/internal/config"
	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

type metricsAgent interface {
	SendMetrics(ctx context.Context, metrics []data.Metrics)
	SendMetric(ctx context.Context, metrics data.Metrics, key string)
}

// Agent struct collects runtume metrics and send them to server
type Agent struct {
	m      *runtime.MemStats
	agent  metricsAgent
	addr   string
	gauges []data.Metrics
}

// NewAgent return Agent structure with endpoint path inside.
func NewAgent(endpoint string) Agent {
	var (
		m  runtime.MemStats
		ms metricsAgent
	)
	cfg := config.GetAgentCfg()
	fmt.Printf("cfg in agent %+v \n", cfg)
	if cfg.GRPC {
		ms = NewGrpcClient(cfg.Adress)
	} else {
		ms = NewHTTPClient("http://" + cfg.Adress)
	}
	a := Agent{m: &m, agent: ms, addr: endpoint, gauges: data.NewGauges()}
	return a
}

// ReadRuntimeMetrics reads runtume metrics.
func (a Agent) ReadRuntimeMetrics(ctx context.Context, interval int) {
	var pollInterval = time.Duration(interval) * time.Second
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("ReadRuntumeMetrics done")
			return
		default:
			runtime.ReadMemStats(a.m)
			time.Sleep(pollInterval)
		}
	}

}

// PullMetrics sends metrics to server.
func (a Agent) PullMetrics(ctx context.Context, interval int, key, crypto string) {

	var reportInterval = time.Duration(interval) * time.Second

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("PullMetrics done")
			return
		default:
			m := a.read()
			a.agent.SendMetrics(ctx, m)
			time.Sleep(reportInterval)
		}
	}
}

func (a Agent) read() []data.Metrics {

	metrics := make([]data.Metrics, 0, 30)

	for _, el := range a.gauges {
		el.ReadValue(a.m)
		metrics = append(metrics, el)
	}

	var (
		d int64   = 1
		f float64 = 0
	)
	metrics = append(metrics, data.Metrics{ID: "PollCount", MType: "counter", Delta: &d, Value: &f})

	return metrics
}
