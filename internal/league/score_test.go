package league

import (
	"path/filepath"
	"testing"

	"townbasketball/internal/domain"
	"townbasketball/internal/store"
)

func TestScoreCorrectionKeepsAudit(t *testing.T) {
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
	if err := l.ScheduleGame(domain.Game{ID: "game-score", HomeTeamID: "team-a", AwayTeamID: "team-b", Venue: "体育馆", Status: "scheduled", Published: true}); err != nil {
		t.Fatal(err)
	}
	updated, _, err := l.PublishScore("game-score", "admin-001", "更正终场比分", 88, 81)
	if err != nil {
		t.Fatal(err)
	}
	if updated.HomeScore != 88 || updated.AwayScore != 81 {
		t.Fatalf("published response retained old score: %+v", updated)
	}
	public, err := l.PublicSchedule()
	if err != nil {
		t.Fatal(err)
	}
	if len(public) != 1 || public[0].HomeScore != 88 || public[0].AwayScore != 81 {
		t.Fatalf("public page retained old score: %+v", public)
	}
	audits, err := l.AuditsForGame("game-score")
	if err != nil {
		t.Fatal(err)
	}
	if len(audits) != 1 || audits[0].NewHome != 88 {
		t.Fatalf("audit missing: %+v", audits)
	}
}
