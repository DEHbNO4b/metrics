// package expert is an adapter between handlers and different stores.
package expert

import (
	"errors"
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

// ExpertConfiguration it is a definition of functional options for configuring Expert struct.
type ExpertConfiguration func(os *Expert) error

// StoreConfig contains main configs of data manegemnt.
type StoreConfig struct {
	Filepath      string
	StoreInterval int
	Restore       bool
}

// Expert struct controlls all work with incoming metrics data.
type Expert struct {
	ram    interfaces.MetricsStorage
	db     interfaces.Database
	config StoreConfig
}

// WithDatabase it as a functional option settings database.
func WithDatabase(db interfaces.Database) ExpertConfiguration {
	return func(e *Expert) error {
		e.db = db
		return nil
	}
}

// WithConfig it as a functional option settings main configs.
func WithConfig(c StoreConfig) ExpertConfiguration {
	return func(e *Expert) error {
		e.config = c
		return nil
	}
}

// WithRam it as a functional option settings RAM store.
func WithRAM(r interfaces.MetricsStorage) ExpertConfiguration {
	return func(e *Expert) error {
		e.ram = r
		return nil
	}
}

// NewExpert return new Expert struct with specifies options.
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

// GetMetrics returns set of metrics from the store.
func (e *Expert) GetMetrics() []data.Metrics {
	return e.ram.GetMetrics()
}

// SetMetric sets the metric value to stor.
func (e *Expert) SetMetric(m data.Metrics) error {
	err := e.setMetricToRAM(m)
	if err != nil {
		return err
	}
	if e.config.StoreInterval == 0 && e.config.Filepath != "" {
		e.StoreData()
	}
	return nil
}

// GetMetric returns specified metric from store.
func (e *Expert) GetMetric(m data.Metrics) (data.Metrics, error) {
	return e.ram.GetMetric(m)
}

// LoadFromDB loads metrics from database to RAM.
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

// StoreData saves metrics from RAM to database.
func (e *Expert) StoreData() {
	e.saveMetricsToDB(e.ram.GetMetrics())
}

func (e *Expert) setMetricToRAM(m data.Metrics) error {
	if m.MType == "counter" {
		c, err := e.ram.GetMetric(m)
		if err != nil && errors.Is(err, interfaces.ErrNotContains) {
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
