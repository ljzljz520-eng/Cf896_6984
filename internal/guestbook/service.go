package guestbook

import (
	"fmt"
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

func (g *Service) Add(nickname, content string) (domain.GuestMessage, error) {
	message := domain.GuestMessage{ID: fmt.Sprintf("msg-%d", g.clock().UnixNano()), Nickname: strings.TrimSpace(nickname), Content: strings.TrimSpace(content), Published: true, CreatedAt: g.clock()}
	if err := message.Validate(); err != nil {
		return domain.GuestMessage{}, err
	}
	if err := g.store.SaveMessage(message); err != nil {
		return domain.GuestMessage{}, err
	}
	return message, nil
}

func (g *Service) ListPublic() ([]domain.GuestMessage, error) {
	messages, err := g.store.ListMessages()
	if err != nil {
		return nil, err
	}
	result := make([]domain.GuestMessage, 0, len(messages))
	for _, message := range messages {
		if message.Published {
			result = append(result, message)
		}
	}
	return result, nil
}

func (g *Service) Moderate(id string, published bool) (domain.GuestMessage, error) {
	message, err := g.store.GetMessage(id)
	if err != nil {
		return domain.GuestMessage{}, err
	}
	message.Published = published
	if err := g.store.SaveMessage(message); err != nil {
		return domain.GuestMessage{}, err
	}
	return message, nil
}

func (g *Service) CountPublic() (int, error) {
	messages, err := g.ListPublic()
	return len(messages), err
}
