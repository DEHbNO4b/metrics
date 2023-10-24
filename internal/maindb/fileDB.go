// package maindb provides functions for storing metrics data.
package maindb

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/DEHbNO4b/metrics/internal/data"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

// FileDB struct implements Database interface.
type FileDB struct {
	filepath string
}

// NewFileDBreturns new FileDB struct.
func NewFileDB(name string) *FileDB {
	return &FileDB{filepath: name}
}

// WriteMetrics it is a method of interfaces.Database. It is writes set of metrics to file.
func (f *FileDB) WriteMetrics(data []data.Metrics) error {
	logger.Log.Info("in fileDB WriteMetrics()")
	file, err := os.OpenFile(filepath.FromSlash(f.filepath), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		logger.Log.Error("unable to open|create storage file ", zap.Error(err))
		return err
	}
	defer file.Close()
	for _, el := range data {
		metric, err := json.Marshal(el)
		if err != nil {
			logger.Log.Error("unable to encode metric ", zap.Error(err))
			return err
		}
		_, err = file.WriteString(string(metric) + "\r\n")
		if err != nil {
			logger.Log.Error("unable to write data to file ", zap.Error(err))
			return err
		}
	}
	return nil
}

// ReadMetrics it is a method of interfaces.Database. It is reads set of metrics from file.
func (f *FileDB) ReadMetrics() ([]data.Metrics, error) {
	logger.Log.Info("in fileDB ReadMetrics()")
	metrics := make([]data.Metrics, 0, 10)
	file, err := os.OpenFile(filepath.FromSlash(f.filepath), os.O_RDONLY, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open storage file, filepath:  ", f.filepath, err.Error())
		return metrics, err
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
		metrics = append(metrics, metric)
	}
	return metrics, nil
}
