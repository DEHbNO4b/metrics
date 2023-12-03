// package server provides management and storing incoming metrics data.
package main

import (
	"os"
	"os/signal"
	"syscall"
	"text/template"

	_ "net/http/pprof"

	"github.com/DEHbNO4b/metrics/internal/app"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

type build struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

const Template = `Build version: {{if .BuildVersion}}{{.buildVersion}}{{else}}N/A{{end}}
Build date: {{if .BuildDate}}{{.buildDate}}{{else}}N/A{{end}}
Build commit: {{if .BuildCommit}}{{.buildCommit}}{{else}}N/A{{end}}
`

func main() {
	b := build{BuildVersion: buildVersion,
		BuildDate:   buildDate,
		BuildCommit: buildCommit}
	t := template.Must(template.New("build").Parse(Template))
	err := t.Execute(os.Stdout, b)
	if err != nil {
		panic(err)
	}
	// TODO:initialize logger
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}
	// TODO: initialize config
	// cfg := config.GetServCfg()
	// TODO: initialize storage
	// withDB, err := selectStore(cfg.Dsn, cfg.StoreFile) //выбор способа храниения данных (sqlDB | fileDB) для эксперта
	// if err != nil {
	// 	panic(err)
	// }
	// expert := expert.NewExpert(expert.WithRAM(maindb.NewMemStorage()), withDB)

	application := app.New()
	go application.Server.MustRun()

	// srv, err := newServer(cfg.Dsn)
	// if err != nil {
	// 	panic(err)
	// }
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	signal := <-stop
	logger.Log.Info("stopped application", zap.Any("signal", signal))
	application.Server.Stop()

	// go func() {
	// 	sigint := make(chan os.Signal, 1)
	// 	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 	<-sigint
	// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()
	// 	if err := srv.Shutdown(ctx); err != nil {
	// 		logger.Log.Error("HTTP Server Shutdown Error", zap.Error(err))
	// 	}
	// 	close(stopped)
	// }()

	// logger.Log.Info("Running server", zap.String("address", cfg.Adress))
	// // start HTTP server
	// if err := srv.ListenAndServe(); err != http.ErrServerClosed {
	// 	logger.Log.Fatal("HTTP server ListenAndServe Error", zap.Error(err))
	// }

	// <-stopped

	logger.Log.Info("Have a nice day!")
}

// func newServer(dsn string) (*http.Server, error) {
// 	cfg := config.GetServCfg()
// 	var ph *handlers.Pinger
// 	if dsn != "" {
// 		sqlDB, err := maindb.NewPostgresDB(dsn)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer sqlDB.Close()
// 		ph = handlers.NewPinger(sqlDB) //хэндлер для пинга
// 	}
// 	withDB, err := selectStore(dsn, cfg.StoreFile) //выбор способа храниения данных (sqlDB | fileDB) для эксперта
// 	if err != nil {
// 		return nil, err
// 	}
// 	expert := expert.NewExpert(expert.WithRAM(maindb.NewMemStorage()), withDB)
// 	defer expert.StoreData() //сохранение данный при завершении программы

// 	h := middleware.Hash{Key: []byte(cfg.HashKey)}
// 	mh := handlers.NewMetrics(expert) //хэндлер для приема и отправки метрик

// 	r := chi.NewRouter()
// 	r.Use(middleware.WithLogging)
// 	r.Use(middleware.WithSubnet)
// 	r.Use(middleware.CryptoHandle)
// 	r.Use(middleware.GzipHandle)
// 	r.Use(h.WithHash)
// 	r.Mount("/debug", middleware.Profiler())
// 	r.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mh.SetMetricsURL))
// 	r.Get(`/value/{type}/{name}`, http.HandlerFunc(mh.GetMetricURL))
// 	r.Post(`/update/`, http.HandlerFunc(mh.SetMetricJSON))
// 	r.Post(`/updates/`, http.HandlerFunc(mh.SetMetricsJSON))
// 	r.Post(`/value/`, http.HandlerFunc(mh.GetMetricJSON))
// 	if dsn != "" {
// 		r.Get(`/ping`, http.HandlerFunc(ph.PingDB))
// 	}
// 	r.Get(`/`, mh.GetMetrics)
// 	srv := &http.Server{
// 		Addr:    cfg.Adress,
// 		Handler: r,
// 	}
// 	return srv, nil
// }
