package maindb

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

type Gauge struct {
	Name string
	Val  float64
}

type Counter struct {
	Name string
	Val  int64
}
type StoreConfig struct {
	Filepath      string
	StoreInterval time.Duration
	Restore       bool
}
type RamStore struct {
	Config   StoreConfig
	Gauges   map[string]float64
	Counters map[string]int64
	DB       interfaces.Database
	sync.RWMutex
}

func NewRamStore(config StoreConfig) *RamStore {
	g := make(map[string]float64)
	c := make(map[string]int64)
	ms := RamStore{Config: config, Gauges: g, Counters: c}
	if config.Restore {
		ms.loadFromStoreFile()
	}
	go ms.storeSchedule()
	return &ms
}
func (ms *RamStore) SetMetric(metric data.Metrics) error {
	switch metric.MType {
	case "gauge":
		ms.Lock()
		ms.Gauges[metric.ID] = *metric.Value
		ms.Unlock()
	case "counter":
		ms.Lock()
		ms.Counters[metric.ID] = ms.Counters[metric.ID] + *metric.Delta
		ms.Unlock()
	default:
		return interfaces.ErrWrongType
	}
	return nil
}
func (ms *RamStore) GetMetrics() []string {
	m := make([]string, 0, 40)
	for name, val := range ms.Gauges {
		m = append(m, name+":"+strconv.FormatFloat(val, 'f', -1, 64))
	}
	for name, val := range ms.Counters {
		m = append(m, name+":"+strconv.FormatInt(val, 10))
	}
	return m
}
func (ms *RamStore) GetMetric(met data.Metrics) (data.Metrics, error) {
	// m := data.Metrics{}
	switch met.MType {
	case "gauge":
		val, ok := ms.Gauges[met.ID]
		if !ok {
			return data.Metrics{}, interfaces.ErrNotContains
		}
		met.Value = &val
	case "counter":
		del, ok := ms.Counters[met.ID]
		if !ok {
			return data.Metrics{}, interfaces.ErrNotContains
		}
		met.Delta = &del
	default:
		return data.Metrics{}, interfaces.ErrWrongType
	}

	return met, nil
}

func (ms *RamStore) GeMetricsData() []data.Metrics {
	metrics := make([]data.Metrics, 0, 30)
	for name, val := range ms.Gauges {
		m := data.NewMetric()
		m.ID = name
		m.MType = "gauge"
		*m.Value = val
		metrics = append(metrics, m)
	}
	for name, val := range ms.Counters {
		m := data.NewMetric()
		m.ID = name
		m.MType = "counter"
		*m.Delta = val
		metrics = append(metrics, m)
	}
	return metrics
}
func (ms *RamStore) loadFromStoreFile() error {
	file, err := os.OpenFile(filepath.FromSlash(ms.Config.Filepath), os.O_RDONLY, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open storage file, filepath:  ", ms.Config.Filepath, err.Error())
		return err
	}
	defer file.Close()
	scaner := bufio.NewScanner(file)
	for scaner.Scan() {
		metric := data.NewMetric()
		line := scaner.Text()
		err := json.Unmarshal([]byte(line), &metric)
		if err != nil {
			logger.Log.Sugar().Error("unable to unmarshal json", err.Error())
			continue
		}
		ms.SetMetric(metric)
	}
	return nil
}
func (ms *RamStore) storeSchedule() {
	for {
		err := ms.StoreData()
		if err != nil {
			return
		}
		time.Sleep(ms.Config.StoreInterval)
	}
}
func (ms *RamStore) StoreData() error {
	file, err := os.OpenFile(filepath.FromSlash(ms.Config.Filepath), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open|create storage file ", err.Error())
		return err
	}
	data := ms.GeMetricsData()
	for _, el := range data {

		metric, err := json.Marshal(el)
		if err != nil {
			logger.Log.Sugar().Error("unable to encode metric ", err.Error())
			continue
		}
		_, err = file.WriteString(string(metric) + "\r\n")
		if err != nil {
			logger.Log.Sugar().Error("unable to write data to file ", err.Error())
			continue
		}
	}
	file.Close()
	return nil
}
