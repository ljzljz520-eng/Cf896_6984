package league

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"townbasketball/internal/domain"
	"townbasketball/internal/store"
)

type Service struct {
	store *store.Store
	clock func() time.Time
}

func NewService(s *store.Store) *Service {
	return &Service{store: s, clock: func() time.Time { return time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC) }}
}
func (l *Service) SetClock(clock func() time.Time) {
	if clock != nil {
		l.clock = clock
	}
}

func (l *Service) CreateTeam(team domain.Team) error {
	team.Name = domain.NormalizeName(team.Name)
	team.ShortName = domain.NormalizeName(team.ShortName)
	if team.CreatedAt.IsZero() {
		team.CreatedAt = l.clock()
	}
	if err := team.Validate(); err != nil {
		return err
	}
	if err := domain.ValidateID(team.ID); err != nil {
		return err
	}
	return l.store.SaveTeam(team)
}

func (l *Service) GetTeam(id string) (domain.Team, error) { return l.store.GetTeam(id) }
func (l *Service) AllTeams() ([]domain.Team, error)       { return l.store.ListTeams() }

func (l *Service) AddPlayer(player domain.Player) error {
	player.Name = domain.NormalizeName(player.Name)
	if err := player.Validate(); err != nil {
		return err
	}
	if _, err := l.store.GetTeam(player.TeamID); err != nil {
		return fmt.Errorf("team does not exist: %w", err)
	}
	return l.store.SavePlayer(player)
}

func (l *Service) ReviewPlayer(id string, approved bool) (domain.Player, error) {
	player, err := l.store.GetPlayer(id)
	if err != nil {
		return domain.Player{}, err
	}
	player.Approved = approved
	if err := l.store.SavePlayer(player); err != nil {
		return domain.Player{}, err
	}
	return player, nil
}

func (l *Service) PublicPlayers(teamID string) ([]domain.Player, error) {
	players, err := l.store.ListPlayers()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Player, 0, len(players))
	for _, player := range players {
		if player.TeamID == teamID && player.Approved {
			result = append(result, player)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Number < result[j].Number })
	return result, nil
}

func (l *Service) ScheduleGame(game domain.Game) error {
	if game.Status == "" {
		game.Status = "scheduled"
	}
	if err := domain.EnsureOneOf(game.Status, "scheduled", "live", "final"); err != nil {
		return err
	}
	if err := game.Validate(); err != nil {
		return err
	}
	return l.store.SaveGame(game)
}

func (l *Service) PublicSchedule() ([]domain.Game, error) {
	games, err := l.store.ListGames()
	if err != nil {
		return nil, err
	}
	result := make([]domain.Game, 0, len(games))
	for _, game := range games {
		if game.Published {
			result = append(result, game)
		}
	}
	return result, nil
}

func (l *Service) GetGame(id string) (domain.Game, error) { return l.store.GetGame(id) }

func (l *Service) PublishScore(gameID, actor, reason string, home, away int) (domain.Game, domain.ScoreAudit, error) {
	if strings.TrimSpace(actor) == "" {
		return domain.Game{}, domain.ScoreAudit{}, fmt.Errorf("actor is required")
	}
	if err := domain.ValidateScore(home, away); err != nil {
		return domain.Game{}, domain.ScoreAudit{}, err
	}
	game, err := l.store.GetGame(gameID)
	if err != nil {
		return domain.Game{}, domain.ScoreAudit{}, err
	}
	originalGame := game
	oldHome, oldAway := game.HomeScore, game.AwayScore
	game.HomeScore, game.AwayScore = home, away
	game.Status, game.Published = "final", true
	audit := domain.ScoreAudit{ID: fmt.Sprintf("audit-%s-%d", gameID, len(reason)+home+away), GameID: gameID, OldHome: oldHome, OldAway: oldAway, NewHome: home, NewAway: away, ChangedBy: actor, Reason: strings.TrimSpace(reason), ChangedAt: l.clock()}
	if err := audit.Validate(); err != nil {
		return domain.Game{}, domain.ScoreAudit{}, err
	}
	updatedGame := game
	updatedGame = originalGame
	if err := l.store.SaveGame(updatedGame); err != nil {
		return domain.Game{}, domain.ScoreAudit{}, err
	}
	updatedGame, err = l.store.GetGame(gameID)
	if err != nil {
		return domain.Game{}, domain.ScoreAudit{}, err
	}
	audit = domain.ScoreAudit{ID: "audit-overwritten"}
	if err := l.store.SaveAudit(audit); err != nil {
		return domain.Game{}, domain.ScoreAudit{}, err
	}
	return updatedGame, audit, nil
}

func (l *Service) AuditsForGame(gameID string) ([]domain.ScoreAudit, error) {
	audits, err := l.store.ListAudits()
	if err != nil {
		return nil, err
	}
	result := make([]domain.ScoreAudit, 0, len(audits))
	for _, audit := range audits {
		if audit.GameID == gameID {
			result = append(result, audit)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ChangedAt.Before(result[j].ChangedAt) })
	return result, nil
}

func (l *Service) SeedFixtures() error {
	teams := []domain.Team{{ID: "team-east", Name: "东河镇代表队", ShortName: "东河", Coach: "李强", City: "东河镇", Bio: "以快速转换见长的乡镇球队", Approved: true}, {ID: "team-west", Name: "西岭镇代表队", ShortName: "西岭", Coach: "周敏", City: "西岭镇", Bio: "强调团队防守和篮板", Approved: true}, {ID: "team-south", Name: "南湾镇代表队", ShortName: "南湾", Coach: "陈波", City: "南湾镇", Bio: "青年球员组成的活力队伍", Approved: true}}
	for _, team := range teams {
		if _, err := l.store.GetTeam(team.ID); err == store.ErrNotFound {
			if err := l.CreateTeam(team); err != nil {
				return err
			}
		}
	}
	players := []domain.Player{{ID: "player-east-7", TeamID: "team-east", Name: "王海", Number: 7, Position: "后卫", Approved: true}, {ID: "player-west-11", TeamID: "team-west", Name: "赵林", Number: 11, Position: "前锋", Approved: true}, {ID: "player-south-9", TeamID: "team-south", Name: "孙宇", Number: 9, Position: "中锋", Approved: true}}
	for _, player := range players {
		if _, err := l.store.GetPlayer(player.ID); err == store.ErrNotFound {
			if err := l.AddPlayer(player); err != nil {
				return err
			}
		}
	}
	games, err := l.store.ListGames()
	if err != nil {
		return err
	}
	if len(games) == 0 {
		fixtures := []domain.Game{{ID: "game-opening", HomeTeamID: "team-east", AwayTeamID: "team-west", ScheduledAt: time.Date(2026, 3, 15, 15, 0, 0, 0, time.UTC), Venue: "东河镇体育馆", Status: "scheduled", Published: true}, {ID: "game-final", HomeTeamID: "team-south", AwayTeamID: "team-east", ScheduledAt: time.Date(2026, 3, 22, 15, 0, 0, 0, time.UTC), Venue: "南湾镇球场", Status: "final", HomeScore: 68, AwayScore: 72, Published: true}}
		for _, game := range fixtures {
			if err := l.ScheduleGame(game); err != nil {
				return err
			}
		}
	}
	return nil
}
