package league

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"townbasketball/internal/domain"
)

type PublicSnapshot struct {
	Teams     []domain.Team       `json:"teams"`
	Games     []domain.Game       `json:"games"`
	Standings []domain.Standing   `json:"standings"`
	Audits    []domain.ScoreAudit `json:"audits"`
}

func (l *Service) QueryGames(filter domain.GameFilter, page domain.Page) ([]domain.Game, error) {
	games, err := l.store.ListGames()
	if err != nil {
		return nil, err
	}
	matched := make([]domain.Game, 0, len(games))
	for _, game := range games {
		if filter.Match(game) {
			matched = append(matched, game)
		}
	}
	return domain.Paginate(matched, page), nil
}

func (l *Service) UpcomingGames(teamID string, now time.Time) ([]domain.Game, error) {
	published := true
	games, err := l.QueryGames(domain.GameFilter{TeamID: teamID, Published: &published, From: now}, domain.Page{Offset: 0, Limit: 100})
	if err != nil {
		return nil, err
	}
	sort.Slice(games, func(i, j int) bool { return games[i].ScheduledAt.Before(games[j].ScheduledAt) })
	return games, nil
}

func (l *Service) RecentResults(teamID string, limit int) ([]domain.Game, error) {
	published := true
	games, err := l.QueryGames(domain.GameFilter{TeamID: teamID, Published: &published, Status: "final"}, domain.Page{Offset: 0, Limit: 100})
	if err != nil {
		return nil, err
	}
	sort.Slice(games, func(i, j int) bool { return games[i].ScheduledAt.After(games[j].ScheduledAt) })
	return domain.Paginate(games, domain.Page{Limit: limit}), nil
}

func (l *Service) PublicSnapshot() (PublicSnapshot, error) {
	teams, err := l.AllTeams()
	if err != nil {
		return PublicSnapshot{}, err
	}
	games, err := l.PublicSchedule()
	if err != nil {
		return PublicSnapshot{}, err
	}
	standings, err := l.Standings()
	if err != nil {
		return PublicSnapshot{}, err
	}
	audits, err := l.store.ListAudits()
	if err != nil {
		return PublicSnapshot{}, err
	}
	publicAudits := make([]domain.ScoreAudit, 0, len(audits))
	for _, audit := range audits {
		if audit.GameID != "" {
			publicAudits = append(publicAudits, domain.PublicAudit(audit))
		}
	}
	return PublicSnapshot{Teams: domain.SortTeamsForDisplay(teams), Games: games, Standings: standings, Audits: publicAudits}, nil
}

func (l *Service) ValidateFixtureSet() error {
	teams, err := l.store.ListTeams()
	if err != nil {
		return err
	}
	games, err := l.store.ListGames()
	if err != nil {
		return err
	}
	if len(teams) < 2 {
		return fmt.Errorf("at least two teams are required")
	}
	if err := domain.ValidateSchedule(games); err != nil {
		return err
	}
	for _, game := range games {
		if _, err := l.store.GetTeam(game.HomeTeamID); err != nil {
			return fmt.Errorf("home team %s: %w", game.HomeTeamID, err)
		}
		if _, err := l.store.GetTeam(game.AwayTeamID); err != nil {
			return fmt.Errorf("away team %s: %w", game.AwayTeamID, err)
		}
	}
	return nil
}

func (l *Service) SearchTeams(query string) ([]domain.Team, error) {
	teams, err := l.store.ListTeams()
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(strings.TrimSpace(query))
	result := make([]domain.Team, 0, len(teams))
	for _, team := range teams {
		if !team.Approved {
			continue
		}
		if q == "" || strings.Contains(strings.ToLower(team.Name), q) || strings.Contains(strings.ToLower(team.City), q) {
			result = append(result, team)
		}
	}
	return domain.SortTeamsForDisplay(result), nil
}
