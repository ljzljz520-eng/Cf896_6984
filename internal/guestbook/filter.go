package guestbook

import (
	"sort"
	"strings"

	"townbasketball/internal/domain"
)

type Filter struct {
	Query     string
	Published *bool
	Limit     int
}

func FilterMessages(messages []domain.GuestMessage, filter Filter) []domain.GuestMessage {
	query := strings.ToLower(strings.TrimSpace(filter.Query))
	result := make([]domain.GuestMessage, 0, len(messages))
	for _, message := range messages {
		if filter.Published != nil && message.Published != *filter.Published {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(message.Content), query) && !strings.Contains(strings.ToLower(message.Nickname), query) {
			continue
		}
		result = append(result, message)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	if filter.Limit > 0 && len(result) > filter.Limit {
		result = result[:filter.Limit]
	}
	return result
}

func GroupByPublication(messages []domain.GuestMessage) map[bool][]domain.GuestMessage {
	grouped := map[bool][]domain.GuestMessage{true: {}, false: {}}
	for _, message := range messages {
		grouped[message.Published] = append(grouped[message.Published], message)
	}
	return grouped
}

func (g *Service) Search(filter Filter) ([]domain.GuestMessage, error) {
	messages, err := g.store.ListMessages()
	if err != nil {
		return nil, err
	}
	return FilterMessages(messages, filter), nil
}
