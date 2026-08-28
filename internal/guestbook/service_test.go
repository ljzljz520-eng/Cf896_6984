package guestbook

import (
	"path/filepath"
	"testing"

	"townbasketball/internal/store"
)

func TestGuestbookModeration(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "league.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	g := NewService(s)
	message, err := g.Add("球迷", "期待精彩比赛")
	if err != nil {
		t.Fatal(err)
	}
	count, err := g.CountPublic()
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal(count)
	}
	if _, err := g.Moderate(message.ID, false); err != nil {
		t.Fatal(err)
	}
	count, err = g.CountPublic()
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal(count)
	}
	if _, err := g.Add("", ""); err == nil {
		t.Fatal("expected validation error")
	}
}
