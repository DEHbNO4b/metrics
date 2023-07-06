package main

import (
	"net/http"

	"github.com/DEHbNO4b/metrics/internal/handlers"
	"github.com/DEHbNO4b/metrics/internal/middlewares"
)

func main() {
	ms := handlers.NewMemStorage()
	serv := http.NewServeMux()
	serv.Handle(`/update/`, middlewares.Conveyor(http.HandlerFunc(ms.SetMetrics), middlewares.IsRightRequest, middlewares.IsPostReq))
	err := http.ListenAndServe(`localhost:8080`, serv)
	if err != nil {
		panic(err)
	}
}
