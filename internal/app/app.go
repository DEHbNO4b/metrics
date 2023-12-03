package app

import (
	appgrpc "github.com/DEHbNO4b/metrics/internal/app/grpc"
	apphttp "github.com/DEHbNO4b/metrics/internal/app/http"
	"github.com/DEHbNO4b/metrics/internal/config"
	"github.com/DEHbNO4b/metrics/internal/expert"
	"github.com/DEHbNO4b/metrics/internal/maindb"
)

type metricServer interface {
	MustRun()
	Stop()
}
type App struct {
	cfg    config.ServerConfig
	Server metricServer
}

func New() *App {
	cfg := config.GetServCfg()
	withDB, err := selectStore(cfg.Dsn, cfg.StoreFile) //выбор способа храниения данных (sqlDB | fileDB) для эксперта
	if err != nil {
		panic(err)
	}
	// TODO: create service layer expert
	expert := expert.NewExpert(expert.WithRAM(maindb.NewMemStorage()), withDB)

	var srv metricServer

	// TODO: create  (gRPC | http)   server
	if cfg.GRPC {
		srv = appgrpc.New(expert, cfg.Adress)
	} else {
		srv = apphttp.New(expert, cfg.Adress)
	}
	return &App{
		cfg:    cfg,
		Server: srv,
	}
}

func selectStore(dsn string, f string) (expert.ExpertConfiguration, error) {
	if dsn != "" {
		db, err := maindb.NewPostgresDB(dsn)
		if err != nil {
			return nil, err
		}
		return expert.WithDatabase(db), nil
	}
	return expert.WithDatabase(maindb.NewFileDB(f)), nil
}
