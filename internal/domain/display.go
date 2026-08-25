package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type GameFilter struct {
	TeamID    string
	Status    string
	Published *bool
	From      time.Time
	To        time.Time
}

type Page struct {
	Offset int
	Limit  int
}

func (p Page) Normalize() Page {
	if p.Offset < 0 {
		p.Offset = 0
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	return p
}

func (f GameFilter) Match(game Game) bool {
	if f.TeamID != "" && game.HomeTeamID != f.TeamID && game.AwayTeamID != f.TeamID {
		return false
	}
	if f.Status != "" && game.Status != f.Status {
		return false
	}
	if f.Published != nil && game.Published != *f.Published {
		return false
	}
	if !f.From.IsZero() && game.ScheduledAt.Before(f.From) {
		return false
	}
	if !f.To.IsZero() && game.ScheduledAt.After(f.To) {
		return false
	}
	return true
}

func Paginate[T any](items []T, page Page) []T {
	page = page.Normalize()
	if page.Offset >= len(items) {
		return []T{}
	}
	end := page.Offset + page.Limit
	if end > len(items) {
		end = len(items)
	}
	return items[page.Offset:end]
}

func TeamDisplayName(team Team) string {
	if team.ShortName != "" {
		return team.ShortName
	}
	return team.Name
}
func GameStatusLabel(status string) string {
	switch status {
	case "scheduled":
		return "未开始"
	case "live":
		return "进行中"
	case "final":
		return "已结束"
	default:
		return "未知状态"
	}
}
func ScoreLabel(game Game) string {
	if game.Status != "final" {
		return "待比赛"
	}
	return fmt.Sprintf("%d : %d", game.HomeScore, game.AwayScore)
}

func SortTeamsForDisplay(teams []Team) []Team {
	copyTeams := append([]Team(nil), teams...)
	sort.SliceStable(copyTeams, func(i, j int) bool {
		a, b := strings.ToLower(TeamDisplayName(copyTeams[i])), strings.ToLower(TeamDisplayName(copyTeams[j]))
		if a != b {
			return a < b
		}
		return copyTeams[i].ID < copyTeams[j].ID
	})
	return copyTeams
}

func SortPlayersForDisplay(players []Player) []Player {
	result := append([]Player(nil), players...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Number != result[j].Number {
			return result[i].Number < result[j].Number
		}
		return result[i].Name < result[j].Name
	})
	return result
}

func PublicTeam(team Team) Team {
	team.Bio = strings.TrimSpace(team.Bio)
	team.CreatedAt = time.Time{}
	return team
}
func PublicPlayer(player Player) Player       { player.SubmittedBy = ""; return player }
func PublicAudit(audit ScoreAudit) ScoreAudit { audit.ChangedBy = ""; return audit }

func ValidateRoster(players []Player) error {
	seenNumbers := make(map[int]struct{}, len(players))
	seenNames := make(map[string]struct{}, len(players))
	for _, player := range players {
		if err := player.Validate(); err != nil {
			return err
		}
		if _, ok := seenNumbers[player.Number]; ok {
			return fmt.Errorf("duplicate jersey number %d", player.Number)
		}
		name := strings.ToLower(strings.TrimSpace(player.Name))
		if _, ok := seenNames[name]; ok {
			return fmt.Errorf("duplicate player name %s", player.Name)
		}
		seenNumbers[player.Number], seenNames[name] = struct{}{}, struct{}{}
	}
	return nil
}

func ValidateSchedule(games []Game) error {
	seen := make(map[string]struct{}, len(games))
	for _, game := range games {
		if _, ok := seen[game.ID]; ok {
			return fmt.Errorf("duplicate game %s", game.ID)
		}
		if err := game.Validate(); err != nil {
			return err
		}
		seen[game.ID] = struct{}{}
	}
	return nil
}
