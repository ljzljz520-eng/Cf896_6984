package auth

import (
	"path/filepath"
	"testing"

	"townbasketball/internal/domain"
	"townbasketball/internal/store"
)

func TestLoginAndRoleChecks(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "league.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	a := NewService(s)
	if err := a.Register(domain.UserAccount{ID: "captain-a", Username: "captain", Role: domain.RoleCaptain, TeamID: "team-a"}, "secret"); err != nil {
		t.Fatal(err)
	}
	session, err := a.Login("captain", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if err := a.RequireCaptain(session, "team-a"); err != nil {
		t.Fatal(err)
	}
	if err := a.RequireCaptain(session, "team-b"); err == nil {
		t.Fatal("expected team mismatch")
	}
	if _, err := a.Login("captain", "wrong"); err == nil {
		t.Fatal("expected login error")
	}
}
