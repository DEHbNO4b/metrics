package appgrpc

import (
	"fmt"
	"net"

	"github.com/DEHbNO4b/metrics/internal/grpc/metrics"
	"github.com/DEHbNO4b/metrics/internal/interfaces"
	logger "github.com/DEHbNO4b/metrics/internal/loger"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	adress     string
}

func New(expert interfaces.MetricsStorage, adress string) *App {
	grpc := grpc.NewServer()
	metrics.Register(grpc, expert)
	return &App{
		gRPCServer: grpc,
		adress:     adress,
	}
}
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func (a *App) Run() error {
	l, err := net.Listen("tcp", a.adress)
	if err != nil {
		return fmt.Errorf("unable to run grpc server %w", err)
	}
	logger.Log.Info("grpc server is running")

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("grpc server error: %w", err)
	}
	return nil
}
func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}
