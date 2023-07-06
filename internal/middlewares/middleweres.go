package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func IsPostReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "wrong method", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func IsRightRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		url := r.URL.Path
		url, _ = strings.CutPrefix(url, "/update/")
		urlValues := strings.Split(url, "/")
		fmt.Println(urlValues)
		fmt.Println(len(urlValues))
		if len(urlValues) < 3 {
			http.Error(w, "bad request", http.StatusNotFound)
			return
		}
		if urlValues[0] == "" {
			http.Error(w, "bad request", http.StatusNotFound)
			return
		}
		if urlValues[1] == "" {
			http.Error(w, "bad request", http.StatusNotFound)
			return
		}
		if urlValues[0] != "counter" && urlValues[0] != "gauge" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
