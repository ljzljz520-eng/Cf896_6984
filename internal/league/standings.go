package league

import (
	"sort"

	"townbasketball/internal/domain"
)

func (l *Service) Standings() ([]domain.Standing, error) {
	teams, err := l.store.ListTeams()
	if err != nil {
		return nil, err
	}
	games, err := l.store.ListGames()
	if err != nil {
		return nil, err
	}
	byID := make(map[string]domain.Standing, len(teams))
	for _, team := range teams {
		if team.Approved {
			byID[team.ID] = domain.Standing{TeamID: team.ID, TeamName: team.Name}
		}
	}
	for _, game := range games {
		if !game.Published || game.Status != "final" {
			continue
		}
		home, okHome := byID[game.HomeTeamID]
		away, okAway := byID[game.AwayTeamID]
		if !okHome || !okAway {
			continue
		}
		home.Played++
		away.Played++
		home.PointsFor += game.HomeScore
		home.PointsAgainst += game.AwayScore
		away.PointsFor += game.AwayScore
		away.PointsAgainst += game.HomeScore
		if game.HomeScore > game.AwayScore {
			home.Wins++
			away.Losses++
			home.TablePoints += 2
		} else if game.HomeScore < game.AwayScore {
			away.Wins++
			home.Losses++
			away.TablePoints += 2
		} else {
			home.TablePoints++
			away.TablePoints++
		}
		byID[game.HomeTeamID], byID[game.AwayTeamID] = home, away
	}
	result := make([]domain.Standing, 0, len(byID))
	for _, standing := range byID {
		result = append(result, standing)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].TablePoints != result[j].TablePoints {
			return result[i].TablePoints > result[j].TablePoints
		}
		diffI := result[i].PointsFor - result[i].PointsAgainst
		diffJ := result[j].PointsFor - result[j].PointsAgainst
		if diffI != diffJ {
			return diffI > diffJ
		}
		return result[i].TeamName < result[j].TeamName
	})
	return result, nil
}

func (l *Service) TeamSummary(teamID string) (domain.Standing, error) {
	standings, err := l.Standings()
	if err != nil {
		return domain.Standing{}, err
	}
	for _, standing := range standings {
		if standing.TeamID == teamID {
			return standing, nil
		}
	}
	return domain.Standing{}, ErrTeamNotFound{ID: teamID}
}

type ErrTeamNotFound struct{ ID string }

func (e ErrTeamNotFound) Error() string { return "team not found: " + e.ID }
