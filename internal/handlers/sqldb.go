package handlers

import (
	"net/http"

	"github.com/DEHbNO4b/metrics/internal/interfaces"
)

type Pinger struct {
	pinger interfaces.Pinger
}

func NewPinger(p interfaces.Pinger) *Pinger {
	return &Pinger{pinger: p}
}
func (p *Pinger) PingDB(w http.ResponseWriter, r *http.Request) {
	if p.pinger == nil {
		http.Error(w, "db is disconected", http.StatusInternalServerError)
		return
	}
	if !p.pinger.Ping() {
		http.Error(w, "db disconected", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
