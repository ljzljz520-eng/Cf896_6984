package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"townbasketball/internal/auth"
	"townbasketball/internal/domain"
	"townbasketball/internal/guestbook"
	"townbasketball/internal/league"
	"townbasketball/internal/media"
	"townbasketball/internal/reporting"
)

type Server struct {
	League    *league.Service
	Auth      *auth.Service
	Media     *media.Service
	Guestbook *guestbook.Service
	Logger    *log.Logger
}

func NewServer(l *league.Service, a *auth.Service, m *media.Service, g *guestbook.Service, logger *log.Logger) *Server {
	return &Server{League: l, Auth: a, Media: m, Guestbook: g, Logger: logger}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.handleHome)
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /readyz", s.handleReadiness)
	mux.HandleFunc("GET /api/public/summary", s.handleSummary)
	mux.HandleFunc("GET /api/public/teams/{id}", s.handleTeam)
	mux.HandleFunc("GET /api/public/messages", s.handleMessages)
	mux.HandleFunc("POST /api/public/messages", s.handleAddMessage)
	mux.HandleFunc("POST /api/captain/roster", s.handleRoster)
	mux.HandleFunc("POST /api/admin/players/{id}/review", s.handleReview)
	mux.HandleFunc("POST /api/admin/games/{id}/score", s.handleScore)
	mux.HandleFunc("GET /api/admin/games/{id}/audits", s.handleAudits)
	return requestLogger(mux, s.Logger)
}

func requestLogger(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if logger != nil {
			logger.Printf("%s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

type publicSummary struct {
	Teams      []domain.Team         `json:"teams"`
	Games      []domain.Game         `json:"games"`
	Standings  []domain.Standing     `json:"standings"`
	Gallery    []domain.MediaAsset   `json:"gallery"`
	Messages   []domain.GuestMessage `json:"messages"`
	GameLabels []string              `json:"game_labels"`
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte("<html><head><title>乡镇篮球联赛</title></head><body><h1>乡镇篮球联赛</h1><p>赛程、球队、比分公告、精彩图片与留言</p><a href=\"/api/public/summary\">查看联赛数据</a></body></html>"))
}

func (s *Server) handleSummary(w http.ResponseWriter, r *http.Request) {
	teams, err := s.League.AllTeams()
	if err != nil {
		writeError(w, 500, err)
		return
	}
	games, err := s.League.PublicSchedule()
	if err != nil {
		writeError(w, 500, err)
		return
	}
	standings, err := s.League.Standings()
	if err != nil {
		writeError(w, 500, err)
		return
	}
	gallery, err := s.Media.PublicGallery()
	if err != nil {
		writeError(w, 500, err)
		return
	}
	messages, err := s.Guestbook.ListPublic()
	if err != nil {
		writeError(w, 500, err)
		return
	}
	teamMap := make(map[string]domain.Team, len(teams))
	for _, team := range teams {
		teamMap[team.ID] = team
	}
	labels := make([]string, 0, len(games))
	for _, game := range games {
		labels = append(labels, reporting.BuildGameLabel(game, teamMap))
	}
	writeJSON(w, 200, publicSummary{Teams: teams, Games: games, Standings: standings, Gallery: gallery, Messages: messages, GameLabels: labels})
}

func (s *Server) handleTeam(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	team, err := s.League.GetTeam(id)
	if err != nil {
		writeError(w, 404, err)
		return
	}
	if !team.Approved {
		writeError(w, 404, fmt.Errorf("team is not public"))
		return
	}
	players, err := s.League.PublicPlayers(id)
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"team": team, "players": players})
}
