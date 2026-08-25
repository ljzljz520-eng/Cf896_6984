package league

import (
	"fmt"
	"sort"
	"time"

	"townbasketball/internal/domain"
)

type HeadToHead struct {
	Games      int `json:"games"`
	HomeWins   int `json:"home_wins"`
	AwayWins   int `json:"away_wins"`
	Draws      int `json:"draws"`
	HomePoints int `json:"home_points"`
	AwayPoints int `json:"away_points"`
}

func CalculateHeadToHead(games []domain.Game, homeID, awayID string) HeadToHead {
	result := HeadToHead{}
	for _, game := range games {
		if !game.Published || game.Status != "final" {
			continue
		}
		if game.HomeTeamID == homeID && game.AwayTeamID == awayID {
			result.Games++
			result.HomePoints += game.HomeScore
			result.AwayPoints += game.AwayScore
			if game.HomeScore > game.AwayScore {
				result.HomeWins++
			} else if game.HomeScore < game.AwayScore {
				result.AwayWins++
			} else {
				result.Draws++
			}
		} else if game.HomeTeamID == awayID && game.AwayTeamID == homeID {
			result.Games++
			result.HomePoints += game.AwayScore
			result.AwayPoints += game.HomeScore
			if game.AwayScore > game.HomeScore {
				result.HomeWins++
			} else if game.AwayScore < game.HomeScore {
				result.AwayWins++
			} else {
				result.Draws++
			}
		}
	}
	return result
}

func WinRate(standing domain.Standing) float64 {
	if standing.Played == 0 {
		return 0
	}
	return float64(standing.Wins) / float64(standing.Played)
}

func VenueUsage(games []domain.Game) map[string]int {
	usage := map[string]int{}
	for _, game := range games {
		venue := game.Venue
		if venue == "" {
			venue = "未指定"
		}
		usage[venue]++
	}
	return usage
}

func ValidateNoOverlap(games []domain.Game) error {
	sorted := append([]domain.Game(nil), games...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ScheduledAt.Before(sorted[j].ScheduledAt) })
	occupied := map[string]time.Time{}
	for _, game := range sorted {
		if previous, ok := occupied[game.Venue]; ok && previous.Equal(game.ScheduledAt) {
			return fmt.Errorf("venue %s has overlapping games", game.Venue)
		}
		occupied[game.Venue] = game.ScheduledAt
	}
	return nil
}

func (l *Service) TeamForm(teamID string) ([]string, error) {
	games, err := l.store.ListGames()
	if err != nil {
		return nil, err
	}
	values := make([]string, 0)
	for _, game := range games {
		if !game.Published || game.Status != "final" {
			continue
		}
		if game.HomeTeamID != teamID && game.AwayTeamID != teamID {
			continue
		}
		if game.HomeTeamID == teamID {
			if game.HomeScore > game.AwayScore {
				values = append(values, "W")
			} else if game.HomeScore < game.AwayScore {
				values = append(values, "L")
			} else {
				values = append(values, "D")
			}
		} else {
			if game.AwayScore > game.HomeScore {
				values = append(values, "W")
			} else if game.AwayScore < game.HomeScore {
				values = append(values, "L")
			} else {
				values = append(values, "D")
			}
		}
	}
	if len(values) > 5 {
		values = values[len(values)-5:]
	}
	return values, nil
}
