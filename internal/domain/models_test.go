package domain

import "testing"

func TestModelsValidateBusinessRules(t *testing.T) {
	if err := (Team{ID: "team-a", Name: "甲队", City: "乡"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Team{ID: "", Name: "甲队", City: "乡"}).Validate(); err == nil {
		t.Fatal("expected team validation error")
	}
	if err := (Player{ID: "p-1", TeamID: "team-a", Name: "球员", Number: 12}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Player{ID: "p-1", TeamID: "team-a", Name: "球员", Number: 100}).Validate(); err == nil {
		t.Fatal("expected number error")
	}
	if err := (Game{ID: "g-1", HomeTeamID: "team-a", AwayTeamID: "team-b", Venue: "场馆"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (Game{ID: "g-1", HomeTeamID: "team-a", AwayTeamID: "team-a", Venue: "场馆"}).Validate(); err == nil {
		t.Fatal("expected same team error")
	}
	if NormalizeName("  东河   镇 ") != "东河 镇" {
		t.Fatal("name normalization failed")
	}
}
