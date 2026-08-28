package httpapi

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := s.Guestbook.ListPublic()
	if err != nil {
		writeError(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"messages": messages})
}

func (s *Server) handleAddMessage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nickname string `json:"nickname"`
		Content  string `json:"content"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, 400, err)
		return
	}
	message, err := s.Guestbook.Add(input.Nickname, input.Content)
	if err != nil {
		writeError(w, 400, err)
		return
	}
	writeJSON(w, 201, message)
}

func parseSession(r *http.Request) (username, password string, err error) {
	username = strings.TrimSpace(r.Header.Get("X-User"))
	password = r.Header.Get("X-Password")
	if username == "" || password == "" {
		return "", "", fmt.Errorf("credentials are required")
	}
	return username, password, nil
}
