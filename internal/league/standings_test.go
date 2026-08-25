package league

import (
	"path/filepath"
	"testing"

	"townbasketball/internal/domain"
	"townbasketball/internal/store"
)

func TestStandingsSortByPointsAndDifference(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "league.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	l := NewService(s)
	for _, team := range []domain.Team{{ID: "team-a", Name: "甲", City: "镇", Approved: true}, {ID: "team-b", Name: "乙", City: "镇", Approved: true}} {
		if err := l.CreateTeam(team); err != nil {
			t.Fatal(err)
		}
	}
	if err := l.ScheduleGame(domain.Game{ID: "game-final", HomeTeamID: "team-a", AwayTeamID: "team-b", Venue: "场", Status: "final", Published: true, HomeScore: 90, AwayScore: 70}); err != nil {
		t.Fatal(err)
	}
	rows, err := l.Standings()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].TeamID != "team-a" || rows[0].Wins != 1 {
		t.Fatalf("unexpected standings: %+v", rows)
	}
}
