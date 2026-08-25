package reporting

import (
	"fmt"
	"strings"

	"townbasketball/internal/domain"
)

type Report struct {
	Title string   `json:"title"`
	Lines []string `json:"lines"`
}

func BuildStandingsReport(rows []domain.Standing) Report {
	lines := make([]string, 0, len(rows))
	for index, row := range rows {
		lines = append(lines, fmt.Sprintf("%d. %s %d分 (%d胜%d负)", index+1, row.TeamName, row.TablePoints, row.Wins, row.Losses))
	}
	return Report{Title: "联赛积分榜", Lines: lines}
}

func BuildGameLabel(game domain.Game, teams map[string]domain.Team) string {
	home, away := teams[game.HomeTeamID], teams[game.AwayTeamID]
	if game.Status == "final" {
		return fmt.Sprintf("%s %d - %d %s", home.ShortName, game.HomeScore, game.AwayScore, away.ShortName)
	}
	return fmt.Sprintf("%s vs %s", home.ShortName, away.ShortName)
}

func ContainsAdminControls(text string) bool {
	return strings.Contains(strings.ToLower(text), "admin") || strings.Contains(text, "后台")
}
