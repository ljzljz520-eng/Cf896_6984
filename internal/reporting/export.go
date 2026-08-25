package reporting

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"townbasketball/internal/domain"
)

type GameRow struct {
	ID     string `json:"id"`
	Match  string `json:"match"`
	Status string `json:"status"`
	Score  string `json:"score"`
}

func BuildGameRows(games []domain.Game, teams map[string]domain.Team) []GameRow {
	rows := make([]GameRow, 0, len(games))
	for _, game := range games {
		rows = append(rows, GameRow{ID: game.ID, Match: BuildGameLabel(game, teams), Status: domain.GameStatusLabel(game.Status), Score: domain.ScoreLabel(game)})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
	return rows
}

func EncodeStandingsJSON(rows []domain.Standing) ([]byte, error) {
	return json.Marshal(BuildStandingsReport(rows))
}

func WriteStandingsCSV(w io.Writer, rows []domain.Standing) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"排名", "球队", "场次", "胜场", "负场", "得分", "失分", "积分"}); err != nil {
		return err
	}
	for index, row := range rows {
		record := []string{fmt.Sprint(index + 1), row.TeamName, fmt.Sprint(row.Played), fmt.Sprint(row.Wins), fmt.Sprint(row.Losses), fmt.Sprint(row.PointsFor), fmt.Sprint(row.PointsAgainst), fmt.Sprint(row.TablePoints)}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func SummarizeTeam(team domain.Team, standing domain.Standing, roster RosterSummaryView) string {
	pieces := []string{team.Name, team.City, fmt.Sprintf("%d场%d胜", standing.Played, standing.Wins), fmt.Sprintf("%d名球员", roster.Approved)}
	return strings.Join(pieces, " · ")
}

type RosterSummaryView struct{ Approved int }

func FilterPublishedMessages(messages []domain.GuestMessage, query string) []domain.GuestMessage {
	q := strings.ToLower(strings.TrimSpace(query))
	result := make([]domain.GuestMessage, 0, len(messages))
	for _, message := range messages {
		if !message.Published {
			continue
		}
		if q == "" || strings.Contains(strings.ToLower(message.Content), q) || strings.Contains(strings.ToLower(message.Nickname), q) {
			result = append(result, message)
		}
	}
	return result
}
