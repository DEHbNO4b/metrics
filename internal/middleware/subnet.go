package middleware

import (
	"net/http"

	"github.com/DEHbNO4b/metrics/internal/config"
)

func WithSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.GetServCfg()
		if cfg.TrustedSubnet != "" {
			ipStr := r.Header.Get("X-Real_IP")
			// ip := net.ParseIP(ipStr)
			// if ip == nil {
			// 	http.Error(w, "wrong ip", http.StatusForbidden)
			// 	return
			// }
			if ipStr != cfg.TrustedSubnet {
				http.Error(w, "wrong ip", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
