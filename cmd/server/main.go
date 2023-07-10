package main

import (
	"net/http"

	"github.com/DEHbNO4b/metrics/internal/data"
	"github.com/DEHbNO4b/metrics/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	parseFlags()
	r := chi.NewRouter()
	ms := data.NewMetStore()
	mh := handlers.NewMetrics(ms)
	// serv := http.NewServeMux()
	// serv.Handle(`/update/`, middlewares.Conveyor(http.HandlerFunc(mh.SetMetrics), middlewares.IsRightRequest, middlewares.IsPostReq))
	r.Post(`/update/{type}/{name}/{value}`, http.HandlerFunc(mh.SetMetrics))
	r.Get(`/value/{type}/{name}`, http.HandlerFunc(mh.GetMetric))
	r.Get(`/`, mh.GetMetrics)
	err := http.ListenAndServe(flagRunAddr, r)
	if err != nil {
		panic(err)
	}
}
