// package server provides management and storing incoming metrics data.
package main

import (
	"net/http"

	_ "net/http/pprof"

	"github.com/DEHbNO4b/metrics/internal/expert"
	"github.com/DEHbNO4b/metrics/internal/handlers"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"github.com/DEHbNO4b/metrics/internal/maindb"
	"github.com/DEHbNO4b/metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	parseFlag()
	sqlDB := maindb.NewPostgresDB(dsn)
	defer sqlDB.Close()
	ph := handlers.NewPinger(sqlDB) //хэндлер для пинга
	config := expert.StoreConfig{
		StoreInterval: storeInterval,
		Filepath:      filestoragepath,
		Restore:       restore,
	}
	withDB := selectStore(dsn, filestoragepath) //выбор способа храниения данных (sqlDB | fileDB) для эксперта
	expert := expert.NewExpert(expert.WithConfig(config), expert.WithRAM(maindb.NewMemStorage()), withDB)
	defer expert.StoreData() //сохранение данный при завершении программы
	h := middleware.Hash{Key: []byte(key)}
	mh := handlers.NewMetrics(expert) //хэндлер для приема и отправки метрик
	r := chi.NewRouter()

	r.Use(middleware.WithLogging)
	r.Use(middleware.GzipHandle)
	r.Use(h.WithHash)
	r.Mount("/debug", middleware.Profiler())
	r.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mh.SetMetricsURL))
	r.Get(`/value/{type}/{name}`, http.HandlerFunc(mh.GetMetricURL))
	r.Post(`/update/`, http.HandlerFunc(mh.SetMetricJSON))
	r.Post(`/updates/`, http.HandlerFunc(mh.SetMetricsJSON))
	r.Post(`/value/`, http.HandlerFunc(mh.GetMetricJSON))
	r.Get(`/ping`, http.HandlerFunc(ph.PingDB))
	r.Get(`/`, mh.GetMetrics)
	logger.Log.Info("Running server", zap.String("address", runAddr))
	err := http.ListenAndServe(runAddr, r)
	if err != nil {
		panic(err)
	}
}

func selectStore(dsn string, f string) expert.ExpertConfiguration {
	if dsn != "" {
		return expert.WithDatabase(maindb.NewPostgresDB(dsn))
	}
	return expert.WithDatabase(maindb.NewFileDB(f))
}
