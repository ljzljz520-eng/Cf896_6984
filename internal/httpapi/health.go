package httpapi

import (
	"net/http"
	"time"
)

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	CheckedAt string `json:"checked_at"`
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok", Service: "town-basketball-league", CheckedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)})
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	if s.League == nil || s.Auth == nil || s.Media == nil || s.Guestbook == nil {
		writeJSON(w, http.StatusServiceUnavailable, healthResponse{Status: "not_ready", Service: "town-basketball-league"})
		return
	}
	writeJSON(w, http.StatusOK, healthResponse{Status: "ready", Service: "town-basketball-league"})
}
