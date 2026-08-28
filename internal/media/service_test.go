package media

import (
	"path/filepath"
	"testing"

	"townbasketball/internal/domain"
	"townbasketball/internal/store"
)

func TestPublicGalleryOnlyPublishedLocalAssets(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "league.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	m := NewService(s)
	if err := m.Add(domain.MediaAsset{ID: "asset-public", Title: "现场", URL: "/media/a.jpg", Published: true}); err != nil {
		t.Fatal(err)
	}
	if err := m.Add(domain.MediaAsset{ID: "asset-draft", Title: "草稿", URL: "/media/b.jpg", Published: false}); err != nil {
		t.Fatal(err)
	}
	assets, err := m.PublicGallery()
	if err != nil {
		t.Fatal(err)
	}
	if len(assets) != 1 || assets[0].ID != "asset-public" {
		t.Fatalf("unexpected assets: %+v", assets)
	}
	if err := m.Add(domain.MediaAsset{ID: "remote", Title: "远端", URL: "https://example.com/x.jpg"}); err == nil {
		t.Fatal("expected local url validation")
	}
}
