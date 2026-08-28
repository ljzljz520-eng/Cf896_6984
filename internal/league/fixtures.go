package league

import (
	"fmt"
	"sort"
	"time"

	"townbasketball/internal/domain"
)

type FixtureRequest struct {
	TeamIDs  []string
	Start    time.Time
	Venue    string
	Interval time.Duration
}

func BuildRoundRobin(request FixtureRequest) ([]domain.Game, error) {
	if len(request.TeamIDs) < 2 {
		return nil, fmt.Errorf("at least two teams are required")
	}
	if request.Start.IsZero() {
		return nil, fmt.Errorf("start time is required")
	}
	if request.Venue == "" {
		return nil, fmt.Errorf("venue is required")
	}
	if request.Interval <= 0 {
		request.Interval = 24 * time.Hour
	}
	ids := append([]string(nil), request.TeamIDs...)
	sort.Strings(ids)
	seen := map[string]struct{}{}
	for _, id := range ids {
		if err := domain.ValidateID(id); err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			return nil, fmt.Errorf("duplicate team %s", id)
		}
		seen[id] = struct{}{}
	}
	games := make([]domain.Game, 0, len(ids)*(len(ids)-1)/2)
	index := 0
	for round := 0; round < len(ids)-1; round++ {
		for offset := round + 1; offset < len(ids); offset++ {
			homeIndex := round
			awayIndex := offset
			if (round+offset)%2 == 1 {
				homeIndex, awayIndex = awayIndex, homeIndex
			}
			game := domain.Game{ID: fmt.Sprintf("game-round-%02d-%02d", round+1, index+1), HomeTeamID: ids[homeIndex], AwayTeamID: ids[awayIndex], ScheduledAt: request.Start.Add(time.Duration(index) * request.Interval), Venue: request.Venue, Status: "scheduled", Published: false}
			games = append(games, game)
			index++
		}
	}
	return games, nil
}

func (l *Service) CreateRoundRobin(request FixtureRequest) ([]domain.Game, error) {
	games, err := BuildRoundRobin(request)
	if err != nil {
		return nil, err
	}
	if err := domain.ValidateSchedule(games); err != nil {
		return nil, err
	}
	if err := l.store.SaveGames(games); err != nil {
		return nil, err
	}
	return games, nil
}

func (l *Service) PublishAnnouncement(announcement domain.Announcement) error {
	announcement.Title = domain.NormalizeName(announcement.Title)
	announcement.Body = domain.NormalizeName(announcement.Body)
	if announcement.CreatedAt.IsZero() {
		announcement.CreatedAt = l.clock()
	}
	if err := announcement.Validate(); err != nil {
		return err
	}
	return l.store.SaveAnnouncement(announcement)
}
func (l *Service) Announcements(publicOnly bool) ([]domain.Announcement, error) {
	values, err := l.store.ListAnnouncements()
	if err != nil {
		return nil, err
	}
	if !publicOnly {
		return values, nil
	}
	result := make([]domain.Announcement, 0, len(values))
	for _, value := range values {
		if value.Published {
			result = append(result, value)
		}
	}
	return result, nil
}
