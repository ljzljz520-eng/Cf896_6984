package media

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

func (m *Service) Add(asset domain.MediaAsset) error {
	asset.Title, asset.Caption = strings.TrimSpace(asset.Title), strings.TrimSpace(asset.Caption)
	if asset.ID == "" || asset.Title == "" || asset.URL == "" {
		return fmt.Errorf("media id, title and url are required")
	}
	if !strings.HasPrefix(asset.URL, "/media/") {
		return fmt.Errorf("media url must be local")
	}
	if asset.CapturedAt.IsZero() {
		asset.CapturedAt = m.clock()
	}
	return m.store.SaveMedia(asset)
}

func (m *Service) Publish(id string, published bool) (domain.MediaAsset, error) {
	asset, err := m.store.GetMedia(id)
	if err != nil {
		return domain.MediaAsset{}, err
	}
	asset.Published = published
	if err := m.store.SaveMedia(asset); err != nil {
		return domain.MediaAsset{}, err
	}
	return asset, nil
}

func (m *Service) PublicGallery() ([]domain.MediaAsset, error) {
	assets, err := m.store.ListMedia()
	if err != nil {
		return nil, err
	}
	result := make([]domain.MediaAsset, 0, len(assets))
	for _, asset := range assets {
		if asset.Published {
			result = append(result, asset)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CapturedAt.After(result[j].CapturedAt) })
	return result, nil
}

func (m *Service) SeedFixtures() error {
	assets, err := m.store.ListMedia()
	if err != nil {
		return err
	}
	if len(assets) > 0 {
		return nil
	}
	fixtures := []domain.MediaAsset{{ID: "media-opening", Title: "开幕式跳球", URL: "/media/opening.jpg", Caption: "联赛开幕式现场", Published: true}, {ID: "media-defense", Title: "防守瞬间", URL: "/media/defense.jpg", Caption: "赛场上的团队协作", Published: true}}
	for _, asset := range fixtures {
		if err := m.Add(asset); err != nil {
			return err
		}
	}
	return nil
}

var _ = store.ErrNotFound
