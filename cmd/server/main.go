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
	// TODO: initialize logger
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}

	// TODO: run application
	application := app.New()
	go application.Server.MustRun()

	// TODO: Gracafull Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	signal := <-stop
	logger.Log.Info("stopped application", zap.Any("signal", signal))
	application.Server.Stop()
	logger.Log.Info("Have a nice day!")
}
