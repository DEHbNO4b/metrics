package maindb

import (
	"os"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

type FileDB struct {
	File *os.File
}

func NewFileDB(name string) *FileDB {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Sugar().Error("unable to open storage file, filepath:  ", ms.Config.Filepath, err.Error())
		return &FileDB{}
	}
}
