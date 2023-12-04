package apphttp

import (
	"context"
	"net/http"
	"time"

	"github.com/DEHbNO4b/metrics/internal/config"
	"github.com/DEHbNO4b/metrics/internal/handlers"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	"github.com/DEHbNO4b/metrics/internal/loger"
	"github.com/DEHbNO4b/metrics/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type App struct {
	cfg    config.ServerConfig
	server *http.Server
	adress string
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func New(expert interfaces.MetricsStorage, adress string) *App {
	cfg := config.GetServCfg()

	h := middleware.Hash{Key: []byte(cfg.HashKey)}
	mh := handlers.NewMetrics(expert) //хэндлер для приема и отправки метрик
	r := chi.NewRouter()
	r.Use(middleware.WithLogging)
	r.Use(middleware.WithSubnet)
	r.Use(middleware.CryptoHandle)
	r.Use(middleware.GzipHandle)
	r.Use(h.WithHash)
	r.Mount("/debug", middleware.Profiler())
	r.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mh.SetMetricsURL))
	r.Get(`/value/{type}/{name}`, http.HandlerFunc(mh.GetMetricURL))
	r.Post(`/update/`, http.HandlerFunc(mh.SetMetricJSON))
	r.Post(`/updates/`, http.HandlerFunc(mh.SetMetricsJSON))
	r.Post(`/value/`, http.HandlerFunc(mh.GetMetricJSON))
	// if cfg.Dsn != "" {
	// 	r.Get(`/ping`, http.HandlerFunc(ph.PingDB))
	// }
	r.Get(`/`, mh.GetMetrics)
	srv := &http.Server{
		Addr:    cfg.Adress,
		Handler: r,
	}

	return &App{
		cfg:    cfg,
		server: srv,
		adress: adress,
	}
}
func (a *App) Run() error {

	loger.Log.Info("Running server", zap.String("address", a.cfg.Adress))

	if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
		loger.Log.Fatal("HTTP server ListenAndServe Error", zap.Error(err))
	}
	return nil
}
func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		loger.Log.Error("HTTP Server Shutdown Error", zap.Error(err))
	}
}
