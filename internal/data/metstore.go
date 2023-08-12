package data

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

var ErrNotContains error = errors.New("not contains this metric")
var ErrWrongType error = errors.New("wrong metric type")

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
type MetStore struct {
	Config   StoreConfig
	Gauges   map[string]float64
	Counters map[string]int64
	sync.RWMutex
}

func NewMetStore(config StoreConfig) *MetStore {
	g := make(map[string]float64)
	c := make(map[string]int64)
	ms := MetStore{Config: config, Gauges: g, Counters: c}
	if config.Restore {
		ms.loadFromStoreFile()
	}
	go ms.storeSchedule()
	return &ms
}
func (ms *MetStore) SetMetric(metric Metrics) error {
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
		return ErrWrongType
	}
	return nil
}
func (ms *MetStore) GetMetrics() []string {
	m := make([]string, 0, 40)
	for name, val := range ms.Gauges {
		m = append(m, name+":"+strconv.FormatFloat(val, 'f', -1, 64))
	}
	for name, val := range ms.Counters {
		m = append(m, name+":"+strconv.FormatInt(val, 10))
	}
	return m
}
func (ms *MetStore) GetMetric(met Metrics) (Metrics, error) {
	// m := Metrics{}
	switch met.MType {
	case "gauge":
		val, ok := ms.Gauges[met.ID]
		if !ok {
			return Metrics{}, ErrNotContains
		}
		met.Value = &val
	case "counter":
		del, ok := ms.Counters[met.ID]
		if !ok {
			return Metrics{}, ErrNotContains
		}
		met.Delta = &del
	default:
		return Metrics{}, ErrWrongType
	}

	return met, nil
}

func (ms *MetStore) GeMetricsData() []Metrics {
	metrics := make([]Metrics, 0, 30)
	for name, val := range ms.Gauges {
		m := NewMetric()
		m.ID = name
		m.MType = "gauge"
		*m.Value = val
		metrics = append(metrics, m)
	}
	for name, val := range ms.Counters {
		m := NewMetric()
		m.ID = name
		m.MType = "counter"
		*m.Delta = val
		metrics = append(metrics, m)
	}
	return metrics
}
func (ms *MetStore) loadFromStoreFile() error {
	file, err := os.OpenFile(filepath.FromSlash(ms.Config.Filepath), os.O_RDONLY, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open storage file, filepath:  ", ms.Config.Filepath, err.Error())
		return err
	}
	defer file.Close()
	scaner := bufio.NewScanner(file)
	for scaner.Scan() {
		metric := NewMetric()
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
func (ms *MetStore) storeSchedule() {
	for {
		err := ms.StoreData()
		if err != nil {
			return
		}
		time.Sleep(ms.Config.StoreInterval)
	}
}
func (ms *MetStore) StoreData() error {
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
