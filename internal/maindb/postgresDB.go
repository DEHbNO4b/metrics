package maindb

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

type MetricsDB struct {
	Config data.StoreConfig
	Db     *sql.DB
}

func NewMetricsDb(dsn string) *MetricsDB {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Log.Error("cannot open db", zap.Error(err))
		return nil
	}
	return &MetricsDB{Db: db}
}
func (mdb *MetricsDB) Ping() error {
	return mdb.Db.Ping()
}
func (mdb *MetricsDB) SetMetric(metric data.Metrics) error {
	switch metric.MType {
	case "gauge":

		// ms.Gauges[metric.ID] = *metric.Value

	case "counter":

		// ms.Counters[metric.ID] = ms.Counters[metric.ID] + *metric.Delta

	default:
		return errors.New("wrong metric type")
	}
	return nil
}
func (mdb *MetricsDB) GetMetrics() []string {
	m := make([]string, 0, 40)
	// for name, val := range ms.Gauges {
	// 	m = append(m, name+":"+strconv.FormatFloat(val, 'f', -1, 64))
	// }
	// for name, val := range ms.Counters {
	// 	m = append(m, name+":"+strconv.FormatInt(val, 10))
	// }
	return m
}
func (mdb *MetricsDB) SetGauge(g data.Gauge) error {
	// ms.Gauges[g.Name] = g.Val
	return nil
}
func (mdb *MetricsDB) SetCounter(c data.Counter) error {
	// ms.Counters[c.Name] = ms.Counters[c.Name] + c.Val
	return nil
}
func (mdb *MetricsDB) GetGauge(name string) (data.Gauge, error) {
	g := data.Gauge{}
	g.Name = name
	// v, ok := ms.Gauges[name]
	// if !ok {
	// 	return g, errors.New("not contains this metric")
	// }
	// g.Val = v
	return g, nil
}
func (mdb *MetricsDB) GetCounter(name string) (data.Counter, error) {
	c := data.Counter{}
	c.Name = name
	// v, ok := ms.Counters[name]
	// if !ok {
	// 	return c, errors.New("not contains this metric")
	// }
	// c.Val = v
	return c, nil
}
func (mdb *MetricsDB) GeMetricsData() []data.Metrics {
	metrics := make([]data.Metrics, 0, 30)
	// for name, val := range ms.Gauges {
	// 	m := NewMetric()
	// 	m.ID = name
	// 	m.MType = "gauge"
	// 	*m.Value = val
	// 	metrics = append(metrics, m)
	// }
	// for name, val := range ms.Counters {
	// 	m := NewMetric()
	// 	m.ID = name
	// 	m.MType = "counter"
	// 	*m.Delta = val
	// 	metrics = append(metrics, m)
	// }
	return metrics
}
func (mdb *MetricsDB) loadFromStoreFile() error {
	file, err := os.OpenFile(filepath.FromSlash(mdb.Config.Filepath), os.O_RDONLY, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open storage file, filepath:  ", mdb.Config.Filepath, err.Error())
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
		mdb.SetMetric(metric)
	}
	return nil
}
func (mdb *MetricsDB) storeSchedule() {
	for {
		err := mdb.StoreData()
		if err != nil {
			return
		}
		time.Sleep(mdb.Config.StoreInterval)
	}
}
func (mdb *MetricsDB) StoreData() error {
	file, err := os.OpenFile(filepath.FromSlash(mdb.Config.Filepath), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open|create storage file ", err.Error())
		return err
	}
	data := mdb.GeMetricsData()
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
