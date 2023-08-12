package maindb

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

type FileDB struct {
	File *os.File
}

func NewFileDB(name string) *FileDB {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log.Error("unable to open|create storage file:  ", zap.Error(err))
		return &FileDB{}
	}
	return &FileDB{File: file}
}
func (f *FileDB) WriteMetrics(data []data.Metrics) error {
	for _, metric := range data {
		err := f.Add(metric)
		if err != nil {
			return err
		}
	}
	return nil
}
func (f *FileDB) ReadMetrics() ([]data.Metrics, error) {
	metrics := make([]data.Metrics, 0, 10)
	scaner := bufio.NewScanner(f.File)
	for scaner.Scan() {
		metric := data.NewMetric()
		line := scaner.Text()
		err := json.Unmarshal([]byte(line), &metric)
		if err != nil {
			logger.Log.Sugar().Error("unable to unmarshal json", err.Error())
			continue
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
func (f *FileDB) Add(m data.Metrics) error {
	metric, err := json.Marshal(m)
	if err != nil {
		logger.Log.Error("unable to encode metric ", zap.Error(err))
		return err
	}
	_, err = f.File.WriteString(string(metric) + "\r\n")
	if err != nil {
		logger.Log.Error("unable to write data to file ", zap.Error(err))
		return err
	}
	return nil
}
