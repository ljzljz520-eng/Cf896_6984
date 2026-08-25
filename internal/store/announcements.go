package store

import (
	"sort"

	"townbasketball/internal/domain"
)

func (s *Store) SaveAnnouncement(announcement domain.Announcement) error {
	return s.save("announcements", announcement.ID, announcement)
}
func (s *Store) GetAnnouncement(id string) (domain.Announcement, error) {
	var value domain.Announcement
	err := s.get("announcements", id, &value)
	return value, err
}
func (s *Store) ListAnnouncements() ([]domain.Announcement, error) {
	var values []domain.Announcement
	err := s.list("announcements", &values)
	sort.Slice(values, func(i, j int) bool { return values[i].CreatedAt.After(values[j].CreatedAt) })
	return values, err
}
