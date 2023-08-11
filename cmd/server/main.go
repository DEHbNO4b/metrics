package main

import (
	"net/http"
	"time"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/handlers"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"github.com/DEHbNO4b/metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {

	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	parseFlag()
	r := chi.NewRouter()
	sc := data.StoreConfig{StoreInterval: time.Duration(storeInterval) * time.Second, Filepath: filestoragepath, Restore: restore}
	ms := data.NewMetStore(sc)
	defer ms.StoreData()
	mh := handlers.NewMetrics(ms)
	r.Use(middleware.WithLogging)
	r.Use(middleware.GzipHandle)
	r.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mh.SetMetricsURL))
	r.Get(`/value/{type}/{name}`, http.HandlerFunc(mh.GetMetricURL))
	r.Post(`/update/`, http.HandlerFunc(mh.SetMetricsJSON))
	r.Post(`/value/`, http.HandlerFunc(mh.GetMetricJSON))
	r.Get(`/`, mh.GetMetrics)
	logger.Log.Info("Running server", zap.String("address", runAddr))
	err := http.ListenAndServe(runAddr, r)
	if err != nil {
		panic(err)
	}
	// router := chi.NewRouter()
	// storeConfig := data.StoreConfig{StoreInterval: time.Duration(storeInterval) * time.Second, Filepath: filestoragepath, Restore: restore}
	// // metricsDB := maindb.NewPostgresDB(dsn)
	// // defer metricsDB.DB.Close()

	// ms := data.NewMetStore(storeConfig)
	// defer ms.StoreData()
	// // if dsn != "" {
	// // 	metricsDB := maindb.NewPostgresDB(dsn)
	// // 	defer metricsDB.DB.Close()
	// // 	store = metricsDB
	// // 	pinger = metricsDB
	// // } else {
	// // 	store = ms
	// // }
	// mhandler := handlers.NewMetrics(ms)
	// // mhandler.Pinger = pinger

	// router.Use(middleware.WithLogging)
	// router.Use(middleware.GzipHandle)
	// router.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mhandler.SetMetricsURL))
	// router.Get(`/value/{type}/{name}`, http.HandlerFunc(mhandler.GetMetricURL))
	// router.Post(`/update/`, http.HandlerFunc(mhandler.SetMetricsJSON))
	// router.Post(`/value/`, http.HandlerFunc(mhandler.GetMetricJSON))
	// router.Get(`/`, mhandler.GetMetrics)
	// router.Get(`/ping`, http.HandlerFunc(mhandler.PingDB))
	// logger.Log.Info("Running server", zap.String("address", runAddr))
	// err := http.ListenAndServe(runAddr, router)
	// if err != nil {
	// 	panic(err)
	// }
}
