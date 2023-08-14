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

type FileDB struct {
	// File *os.File
	filepath string
}

func NewFileDB(name string) *FileDB {
	// file, err := os.OpenFile(filepath.FromSlash(name), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	logger.Log.Error("unable to open|create storage file:  ", zap.Error(err))
	// 	panic(err)
	// 	// return nil
	// }
	return &FileDB{filepath: name}
}
func (f *FileDB) WriteMetrics(data []data.Metrics) error {
	file, err := os.OpenFile(filepath.FromSlash(f.filepath), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open|create storage file ", err.Error())
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

func (f *FileDB) ReadMetrics() ([]data.Metrics, error) {
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
func (f *FileDB) Add(m data.Metrics) error {
	// metric, err := json.Marshal(m)
	// if err != nil {
	// 	logger.Log.Error("unable to encode metric ", zap.Error(err))
	// 	return err
	// }
	// _, err = f.File.WriteString(string(metric) + "\r\n")
	// if err != nil {
	// 	logger.Log.Error("unable to write data to file ", zap.Error(err))
	// 	return err
	// }
	return nil
}
