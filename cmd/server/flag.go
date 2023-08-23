package main

import (
	"flag"
	"os"
	"strconv"

	logger "github.com/DEHbNO4b/metrics/internal/loger"
)

var (
	runAddr         string
	filestoragepath string
	dsn             string
	key             string
	storeInterval   int
	restore         bool
)

func parseFlag() {
	flag.StringVar(&runAddr, "a", "localhost:8080", "adress and port for running")
	flag.StringVar(&key, "k", "", "hash key")
	flag.StringVar(&filestoragepath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.StringVar(&dsn, "d", "", "dsn for postgres")
	flag.IntVar(&storeInterval, "i", 300, "data store interval")
	flag.BoolVar(&restore, "r", true, "restore_flag")
	flag.Parse()
	if ep := os.Getenv("ADDRESS"); ep != "" {
		runAddr = ep
	}
	if si := os.Getenv("STORE_INTERVAL"); si != "" {
		sInt, err := strconv.Atoi(si)
		if err != nil {
			logger.Log.Sugar().Error("unable to convert STORE_INTERVAL to int", err.Error())
			return
		}
		storeInterval = sInt
	}
	if fp := os.Getenv("FILE_STORAGE_PATH"); fp != "" {
		filestoragepath = fp
	}
	if k := os.Getenv("KEY"); k != "" {
		key = k
	}
	if dbdsn := os.Getenv("DATABASE_DSN"); dbdsn != "" {
		dsn = dbdsn
	}
	if r := os.Getenv("RESTORE"); r != "" {
		re, err := strconv.ParseBool(r)
		if err != nil {
			logger.Log.Sugar().Error("unable to convert STORE_INTERVAL to int", err.Error())
			return
		}
		restore = re
	}
}
