package league

import (
	"path/filepath"
	"testing"

	"townbasketball/internal/domain"
	"townbasketball/internal/store"
)

func newLeagueTest(t *testing.T) (*Service, *store.Store) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "league.db"))
	if err != nil {
		t.Fatal(err)
	}
	return NewService(s), s
}

func TestRosterReviewWorkflow(t *testing.T) {
	l, s := newLeagueTest(t)
	defer s.Close()
	if err := l.CreateTeam(domain.Team{ID: "team-a", Name: "甲队", City: "镇", Approved: true}); err != nil {
		t.Fatal(err)
	}
	player := domain.Player{ID: "player-a", TeamID: "team-a", Name: "新球员", Number: 8, Position: "后卫"}
	if err := l.AddPlayer(player); err != nil {
		t.Fatal(err)
	}
	public, err := l.PublicPlayers("team-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(public) != 0 {
		t.Fatal("unreviewed player is public")
	}
	if _, err := l.ReviewPlayer(player.ID, true); err != nil {
		t.Fatal(err)
	}
	public, err = l.PublicPlayers("team-a")
	if err != nil {
		t.Fatal(err)
	}
	if len(public) != 1 || !public[0].Approved {
		t.Fatalf("unexpected players: %+v", public)
	}
}

func TestScheduleFiltersUnpublishedGames(t *testing.T) {
	l, s := newLeagueTest(t)
	defer s.Close()
	for _, team := range []domain.Team{{ID: "team-a", Name: "甲", City: "镇", Approved: true}, {ID: "team-b", Name: "乙", City: "镇", Approved: true}} {
		if err := l.CreateTeam(team); err != nil {
			t.Fatal(err)
		}
	}
	game := domain.Game{ID: "game-hidden", HomeTeamID: "team-a", AwayTeamID: "team-b", Venue: "体育馆", Published: false}
	if err := l.ScheduleGame(game); err != nil {
		t.Fatal(err)
	}
	games, err := l.PublicSchedule()
	if err != nil {
		t.Fatal(err)
	}
	if len(games) != 0 {
		t.Fatal("unpublished game leaked")
	}
}
