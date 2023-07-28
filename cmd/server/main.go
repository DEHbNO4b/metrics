package main

import (
	"net/http"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/handlers"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"github.com/DEHbNO4b/metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	parseFlag()
	r := chi.NewRouter()
	mh := handlers.NewMetrics(data.NewMetStore())
	r.Use(middleware.WithLogging)
	r.Use(middleware.GzipHandle)
	r.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mh.SetMetricsURL))
	r.Post(`/update/`, http.HandlerFunc(mh.SetMetricsJSON))
	r.Post(`/value/`, http.HandlerFunc(mh.GetMetricJSON))
	r.Get(`/value/{type}/{name}`, http.HandlerFunc(mh.GetMetricURL))
	r.Get(`/`, mh.GetMetrics)
	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
