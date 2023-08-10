package main

import (
	"net/http"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/handlers"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"github.com/DEHbNO4b/metrics/internal/maindb"
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
	router := chi.NewRouter()
	// storeConfig := data.StoreConfig{StoreInterval: time.Duration(storeInterval) * time.Second, Filepath: filestoragepath, Restore: restore}
	postgresDB := maindb.NewPostgresDB(dsn)
	defer postgresDB.DB.Close()
	ms := data.NewMetStore(storeConfig)
	defer ms.StoreData()
	mhandler := handlers.NewMetrics(ms)
	mhandler.Pinger = postgresDB
	router.Use(middleware.WithLogging)
	router.Use(middleware.GzipHandle)
	router.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mhandler.SetMetricsURL))
	router.Get(`/value/{type}/{name}`, http.HandlerFunc(mhandler.GetMetricURL))
	router.Post(`/update/`, http.HandlerFunc(mhandler.SetMetricsJSON))
	router.Post(`/value/`, http.HandlerFunc(mhandler.GetMetricJSON))
	router.Get(`/`, mhandler.GetMetrics)
	router.Get(`/ping`, http.HandlerFunc(mhandler.PingDB))
	logger.Log.Info("Running server", zap.String("address", runAddr))
	err := http.ListenAndServe(runAddr, router)
	if err != nil {
		panic(err)
	}
}
