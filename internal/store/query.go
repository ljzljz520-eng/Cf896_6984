package store

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
	"townbasketball/internal/domain"
)

type Snapshot struct {
	Teams    []domain.Team         `json:"teams"`
	Players  []domain.Player       `json:"players"`
	Games    []domain.Game         `json:"games"`
	Audits   []domain.ScoreAudit   `json:"audits"`
	Media    []domain.MediaAsset   `json:"media"`
	Messages []domain.GuestMessage `json:"messages"`
}

func (s *Store) Count(bucket string) (int, error) {
	count := 0
	err := s.transaction(false, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		count = b.Stats().KeyN
		return nil
	})
	return count, err
}

func (s *Store) Keys(bucket, prefix string) ([]string, error) {
	keys := []string{}
	err := s.transaction(false, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		return b.ForEach(func(k, v []byte) error {
			if v != nil && strings.HasPrefix(string(k), prefix) {
				keys = append(keys, string(k))
			}
			return nil
		})
	})
	return keys, err
}

func (s *Store) SaveTeams(teams []domain.Team) error {
	return saveBatch(s, "teams", teams, func(v domain.Team) string { return v.ID })
}
func (s *Store) SavePlayers(players []domain.Player) error {
	return saveBatch(s, "players", players, func(v domain.Player) string { return v.ID })
}
func (s *Store) SaveGames(games []domain.Game) error {
	return saveBatch(s, "games", games, func(v domain.Game) string { return v.ID })
}

func saveBatch[T any](s *Store, bucket string, values []T, id func(T) string) error {
	return s.transaction(true, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		for _, value := range values {
			keyValue := id(value)
			if strings.TrimSpace(keyValue) == "" {
				return fmt.Errorf("key required")
			}
			data, err := json.Marshal(value)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(keyValue), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Snapshot() (Snapshot, error) {
	var snapshot Snapshot
	if err := s.list("teams", &snapshot.Teams); err != nil {
		return snapshot, err
	}
	if err := s.list("players", &snapshot.Players); err != nil {
		return snapshot, err
	}
	if err := s.list("games", &snapshot.Games); err != nil {
		return snapshot, err
	}
	if err := s.list("audits", &snapshot.Audits); err != nil {
		return snapshot, err
	}
	if err := s.list("media", &snapshot.Media); err != nil {
		return snapshot, err
	}
	if err := s.list("messages", &snapshot.Messages); err != nil {
		return snapshot, err
	}
	return snapshot, nil
}

func (s *Store) ExportJSON() ([]byte, error) {
	snapshot, err := s.Snapshot()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(snapshot, "", "  ")
}

func (s *Store) SetSetting(name, value string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("setting name is required")
	}
	return s.save("settings", name, map[string]string{"name": name, "value": value, "updated_at": time.Now().UTC().Format(time.RFC3339)})
}
func (s *Store) GetSetting(name string) (string, error) {
	var value struct {
		Value string `json:"value"`
	}
	if err := s.get("settings", name, &value); err != nil {
		return "", err
	}
	return value.Value, nil
}
