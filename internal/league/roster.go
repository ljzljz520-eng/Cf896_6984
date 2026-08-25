package league

import (
	"fmt"
	"sort"
	"strings"

	"townbasketball/internal/domain"
)

type RosterSummary struct {
	TeamID    string         `json:"team_id"`
	Total     int            `json:"total"`
	Approved  int            `json:"approved"`
	Pending   int            `json:"pending"`
	Positions map[string]int `json:"positions"`
}

func (l *Service) SubmitRoster(teamID, submitter string, players []domain.Player) ([]domain.Player, error) {
	if strings.TrimSpace(teamID) == "" || strings.TrimSpace(submitter) == "" {
		return nil, fmt.Errorf("team and submitter are required")
	}
	if err := domain.ValidateRoster(players); err != nil {
		return nil, err
	}
	for i := range players {
		players[i].TeamID = teamID
		players[i].SubmittedBy = submitter
		players[i].Approved = false
		if err := l.AddPlayer(players[i]); err != nil {
			return nil, err
		}
	}
	return players, nil
}

func (l *Service) ReviewRoster(teamID string, approved bool) ([]domain.Player, error) {
	players, err := l.store.ListPlayers()
	if err != nil {
		return nil, err
	}
	reviewed := make([]domain.Player, 0)
	for _, player := range players {
		if player.TeamID == teamID {
			player.Approved = approved
			if err := l.store.SavePlayer(player); err != nil {
				return nil, err
			}
			reviewed = append(reviewed, player)
		}
	}
	sort.Slice(reviewed, func(i, j int) bool { return reviewed[i].Number < reviewed[j].Number })
	return reviewed, nil
}

func (l *Service) RosterSummary(teamID string) (RosterSummary, error) {
	players, err := l.store.ListPlayers()
	if err != nil {
		return RosterSummary{}, err
	}
	summary := RosterSummary{TeamID: teamID, Positions: map[string]int{}}
	for _, player := range players {
		if player.TeamID != teamID {
			continue
		}
		summary.Total++
		if player.Approved {
			summary.Approved++
		} else {
			summary.Pending++
		}
		position := strings.TrimSpace(player.Position)
		if position == "" {
			position = "未分配"
		}
		summary.Positions[position]++
	}
	return summary, nil
}

func (l *Service) RemovePlayer(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("player id is required")
	}
	return l.store.Delete("players", id)
}

func (l *Service) TeamProfile(teamID string) (domain.Team, []domain.Player, RosterSummary, error) {
	team, err := l.GetTeam(teamID)
	if err != nil {
		return domain.Team{}, nil, RosterSummary{}, err
	}
	players, err := l.PublicPlayers(teamID)
	if err != nil {
		return domain.Team{}, nil, RosterSummary{}, err
	}
	summary, err := l.RosterSummary(teamID)
	if err != nil {
		return domain.Team{}, nil, RosterSummary{}, err
	}
	return domain.PublicTeam(team), players, summary, nil
}
