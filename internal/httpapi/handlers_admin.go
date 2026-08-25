package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"townbasketball/internal/auth"
	"townbasketball/internal/domain"
)

func (s *Server) authenticate(r *http.Request) (domain.UserAccount, error) {
	username, password, err := parseSession(r)
	if err != nil {
		return domain.UserAccount{}, err
	}
	session, err := s.Auth.Login(username, password)
	if err != nil {
		return domain.UserAccount{}, err
	}
	return domain.UserAccount{ID: session.UserID, Username: session.Username, Role: session.Role, TeamID: session.TeamID, Active: true}, nil
}

func (s *Server) handleRoster(w http.ResponseWriter, r *http.Request) {
	user, err := s.authenticate(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	var input struct {
		ID       string `json:"id"`
		TeamID   string `json:"team_id"`
		Name     string `json:"name"`
		Number   int    `json:"number"`
		Position string `json:"position"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, 400, err)
		return
	}
	if err := s.Auth.RequireCaptain(authSession(user), input.TeamID); err != nil {
		writeError(w, 403, err)
		return
	}
	player := domain.Player{ID: strings.TrimSpace(input.ID), TeamID: input.TeamID, Name: input.Name, Number: input.Number, Position: input.Position, SubmittedBy: user.ID}
	if err := s.League.AddPlayer(player); err != nil {
		writeError(w, 400, err)
		return
	}
	writeJSON(w, 201, player)
}

func (s *Server) handleReview(w http.ResponseWriter, r *http.Request) {
	user, err := s.authenticate(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	if err := s.Auth.RequireAdmin(authSession(user)); err != nil {
		writeError(w, 403, err)
		return
	}
	approved, err := strconv.ParseBool(r.URL.Query().Get("approved"))
	if err != nil {
		writeError(w, 400, fmt.Errorf("approved must be boolean"))
		return
	}
	player, err := s.League.ReviewPlayer(r.PathValue("id"), approved)
	if err != nil {
		writeError(w, 404, err)
		return
	}
	writeJSON(w, 200, player)
}

func (s *Server) handleScore(w http.ResponseWriter, r *http.Request) {
	user, err := s.authenticate(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	if err := s.Auth.RequireAdmin(authSession(user)); err != nil {
		writeError(w, 403, err)
		return
	}
	var input struct {
		Home   int    `json:"home"`
		Away   int    `json:"away"`
		Reason string `json:"reason"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, 400, err)
		return
	}
	game, audit, err := s.League.PublishScore(r.PathValue("id"), user.ID, input.Reason, input.Home, input.Away)
	if err != nil {
		writeError(w, 400, err)
		return
	}
	writeJSON(w, 200, map[string]any{"game": game, "audit": audit})
}

func (s *Server) handleAudits(w http.ResponseWriter, r *http.Request) {
	user, err := s.authenticate(r)
	if err != nil {
		writeError(w, 401, err)
		return
	}
	if err := s.Auth.RequireAdmin(authSession(user)); err != nil {
		writeError(w, 403, err)
		return
	}
	audits, err := s.League.AuditsForGame(r.PathValue("id"))
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"audits": audits})
}

func authSession(user domain.UserAccount) auth.Session {
	return auth.Session{UserID: user.ID, Username: user.Username, Role: user.Role, TeamID: user.TeamID}
}
