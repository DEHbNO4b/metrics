package main

import (
	"net/http"
	"time"

	"github.com/DEHbNO4b/metrics/internal/handlers"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"github.com/DEHbNO4b/metrics/internal/maindb"
	"github.com/DEHbNO4b/metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var db interfaces.Database

func main() {
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	parseFlag()
	sqlDB := maindb.NewPostgresDB(dsn)
	defer sqlDB.DB.Close()
	filedb := maindb.NewFileDB(filestoragepath)
	defer filedb.File.Close()
	sc := maindb.StoreConfig{
		StoreInterval: time.Duration(storeInterval) * time.Second,
		Filepath:      filestoragepath,
		Restore:       restore,
	}
	if sqlDB.Ping() == nil {
		db = sqlDB
	} else {
		db = filedb
	}
	rs := maindb.NewRAMStore(sc, db)
	defer rs.StoreData()
	mh := handlers.NewMetrics(rs)
	mh.Pinger = sqlDB
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.GzipHandle)
	r.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mh.SetMetricsURL))
	r.Get(`/value/{type}/{name}`, http.HandlerFunc(mh.GetMetricURL))
	r.Post(`/update/`, http.HandlerFunc(mh.SetMetricJSON))
	r.Post(`/updates/`, http.HandlerFunc(mh.SetMetricsJSON))
	r.Post(`/value/`, http.HandlerFunc(mh.GetMetricJSON))
	r.Get(`/ping`, http.HandlerFunc(mh.PingDB))
	r.Get(`/`, mh.GetMetrics)
	logger.Log.Info("Running server", zap.String("address", runAddr))
	err := http.ListenAndServe(runAddr, r)
	if err != nil {
		panic(err)
	}
}
