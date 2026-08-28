package httpapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"townbasketball/internal/auth"
	"townbasketball/internal/guestbook"
	"townbasketball/internal/league"
	"townbasketball/internal/media"
	"townbasketball/internal/store"
)

func newHTTPTest(t *testing.T) (*Server, *store.Store) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "league.db"))
	if err != nil {
		t.Fatal(err)
	}
	a := auth.NewService(s)
	l := league.NewService(s)
	m := media.NewService(s)
	g := guestbook.NewService(s)
	return NewServer(l, a, m, g, log.New(io.Discard, "", 0)), s
}

func TestWorkflowPublicScoreboard(t *testing.T) {
	server, s := newHTTPTest(t)
	defer s.Close()
	l := server.League
	if err := l.SeedFixtures(); err != nil {
		t.Fatal(err)
	}
	if err := server.Media.SeedFixtures(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/public/summary", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatal(rec.Code)
	}
	var body struct {
		Games   []any `json:"games"`
		Gallery []any `json:"gallery"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Games) != 2 || len(body.Gallery) != 2 {
		t.Fatalf("unexpected summary: %+v", body)
	}
	if strings.Contains(strings.ToLower(rec.Body.String()), "admin") {
		t.Fatal("public payload leaked admin controls")
	}
}

func TestWorkflowRosterReview(t *testing.T) {
	server, s := newHTTPTest(t)
	defer s.Close()
	if err := server.League.SeedFixtures(); err != nil {
		t.Fatal(err)
	}
	if err := server.Auth.SeedDefaults(); err != nil {
		t.Fatal(err)
	}
	payload := `{"id":"player-http","team_id":"team-east","name":"网页球员","number":18,"position":"后卫"}`
	req := httptest.NewRequest(http.MethodPost, "/api/captain/roster", strings.NewReader(payload))
	req.Header.Set("X-User", "captain-east")
	req.Header.Set("X-Password", "captain-pass")
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("submit status %d body %s", rec.Code, rec.Body.String())
	}
	req = httptest.NewRequest(http.MethodGet, "/api/public/teams/team-east", nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if strings.Contains(rec.Body.String(), "网页球员") {
		t.Fatal("unapproved roster leaked")
	}
}

func TestWorkflowGuestbook(t *testing.T) {
	server, s := newHTTPTest(t)
	defer s.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/public/messages", strings.NewReader(`{"nickname":"观众","content":"加油"}`))
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatal(rec.Code)
	}
	req = httptest.NewRequest(http.MethodGet, "/api/public/messages", nil)
	rec = httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)
	if !strings.Contains(rec.Body.String(), "观众") {
		t.Fatal("message not listed")
	}
}
