package expert

import (
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

type ExpertConfiguration func(os *Expert) error

type StoreConfig struct {
	Filepath      string
	StoreInterval int
	Restore       bool
}
type Expert struct {
	ram    interfaces.MetricsStorage
	db     interfaces.Database
	config StoreConfig
}

func WithDatabase(db interfaces.Database) ExpertConfiguration {
	return func(e *Expert) error {
		e.db = db
		return nil
	}
}
func WithConfig(c StoreConfig) ExpertConfiguration {
	return func(e *Expert) error {
		e.config = c
		return nil
	}
}
func WithRAM(r interfaces.MetricsStorage) ExpertConfiguration {
	return func(e *Expert) error {
		e.ram = r
		return nil
	}
}

func NewExpert(cfgs ...ExpertConfiguration) *Expert {
	e := &Expert{}
	for _, cfg := range cfgs {
		err := cfg(e)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
	if e.config.Restore && e.config.Filepath != "" {
		e.LoadFromDB()
	}
	if e.config.StoreInterval > 0 && e.config.Filepath != "" {
		go func() {
			for {
				e.StoreData()
				time.Sleep(time.Duration(e.config.StoreInterval) * time.Second)
			}
		}()
	}
	return e
}

func (e *Expert) GetMetrics() []data.Metrics {
	return e.ram.GetMetrics()
}
func (e *Expert) SetMetric(m data.Metrics) error {
	e.setMetricToRAM(m)
	if e.config.StoreInterval == 0 && e.config.Filepath != "" {
		// e.saveMetricsToDB(e.ram.GetMetrics())
		e.StoreData()
	}
	return nil
}
func (e *Expert) GetMetric(m data.Metrics) (data.Metrics, error) {
	return e.ram.GetMetric(m)
}
func (e *Expert) LoadFromDB() error {
	metrics, err := e.db.ReadMetrics()
	if err != nil {
		return err
	}
	for _, el := range metrics {
		e.ram.SetMetric(el)
	}
	return nil
}
func (e *Expert) setMetricToRAM(m data.Metrics) error {
	if m.MType == "counter" {
		c, err := e.ram.GetMetric(m)
		if err != nil && err == interfaces.ErrNotContains {
			return e.ram.SetMetric(m)
		}
		*m.Delta = *c.Delta + *m.Delta
		return e.ram.SetMetric(m)
	}
	return e.ram.SetMetric(m)
}
func (e *Expert) saveMetricsToDB(m []data.Metrics) {
	e.db.WriteMetrics(m)
}
func (e *Expert) StoreData() {
	e.saveMetricsToDB(e.ram.GetMetrics())
}
