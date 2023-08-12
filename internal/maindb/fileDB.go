package maindb

import (
	"os"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"go.uber.org/zap"
)

type FileDB struct {
	File *os.File
}

func NewFileDB(name string) *FileDB {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("unable to open storage file, filepath:  ", zap.Error(err))
		return &FileDB{}
	}
	return &FileDB{File: file}
}
