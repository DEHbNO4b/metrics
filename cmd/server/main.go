package main

import (
	"net/http"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/handlers"
	"github.com/DEHbNO4b/metrics/internal/middlewares"
)

func main() {
	ms := data.NewMetStore()
	mh := handlers.NewMetrics(ms)
	serv := http.NewServeMux()
	serv.Handle(`/update/`, middlewares.Conveyor(http.HandlerFunc(mh.SetMetrics), middlewares.IsRightRequest, middlewares.IsPostReq))
	err := http.ListenAndServe(`localhost:8080`, serv)
	if err != nil {
		panic(err)
	}
}
