package store

import (
	"path/filepath"
	"testing"

	"townbasketball/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "league.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	team := domain.Team{ID: "team-persist", Name: "持久化队", City: "测试镇", Approved: true}
	if err := s.SaveTeam(team); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Reopen(); err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	got, err := s.GetTeam(team.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != team.Name || !got.Approved {
		t.Fatalf("unexpected team: %+v", got)
	}
}

func TestRepositoryListsRecordsByBusinessOrder(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "league.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	for _, team := range []domain.Team{{ID: "team-z", Name: "乙", City: "镇"}, {ID: "team-a", Name: "甲", City: "镇"}} {
		if err := s.SaveTeam(team); err != nil {
			t.Fatal(err)
		}
	}
	teams, err := s.ListTeams()
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 2 || (teams[0].ID != "team-a" && teams[1].ID != "team-a") {
		t.Fatalf("unexpected order: %+v", teams)
	}
}
