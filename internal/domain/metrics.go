package domain

import (
	"fmt"
	"sort"
	"strings"
)

type TeamMetrics struct {
	TeamID         string   `json:"team_id"`
	Played         int      `json:"played"`
	AverageFor     float64  `json:"average_for"`
	AverageAgainst float64  `json:"average_against"`
	Form           []string `json:"form"`
}

func CalculateTeamMetrics(teamID string, games []Game) TeamMetrics {
	metrics := TeamMetrics{TeamID: teamID}
	relevant := make([]Game, 0)
	for _, game := range games {
		if !game.Published || game.Status != "final" {
			continue
		}
		if game.HomeTeamID == teamID || game.AwayTeamID == teamID {
			relevant = append(relevant, game)
		}
	}
	sort.Slice(relevant, func(i, j int) bool { return relevant[i].ScheduledAt.Before(relevant[j].ScheduledAt) })
	for _, game := range relevant {
		metrics.Played++
		if game.HomeTeamID == teamID {
			metrics.AverageFor += float64(game.HomeScore)
			metrics.AverageAgainst += float64(game.AwayScore)
			if game.HomeScore > game.AwayScore {
				metrics.Form = append(metrics.Form, "W")
			} else if game.HomeScore < game.AwayScore {
				metrics.Form = append(metrics.Form, "L")
			} else {
				metrics.Form = append(metrics.Form, "D")
			}
		} else {
			metrics.AverageFor += float64(game.AwayScore)
			metrics.AverageAgainst += float64(game.HomeScore)
			if game.AwayScore > game.HomeScore {
				metrics.Form = append(metrics.Form, "W")
			} else if game.AwayScore < game.HomeScore {
				metrics.Form = append(metrics.Form, "L")
			} else {
				metrics.Form = append(metrics.Form, "D")
			}
		}
	}
	if metrics.Played > 0 {
		metrics.AverageFor /= float64(metrics.Played)
		metrics.AverageAgainst /= float64(metrics.Played)
	}
	if len(metrics.Form) > 5 {
		metrics.Form = metrics.Form[len(metrics.Form)-5:]
	}
	return metrics
}

func RenderMetrics(metrics TeamMetrics) string {
	form := strings.Join(metrics.Form, "")
	return fmt.Sprintf("%s %.1f/%.1f %s", metrics.TeamID, metrics.AverageFor, metrics.AverageAgainst, form)
}
