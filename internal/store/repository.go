package store

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	bolt "go.etcd.io/bbolt"
	"townbasketball/internal/domain"
)

func (s *Store) SaveTeam(team domain.Team) error { return s.save("teams", team.ID, team) }
func (s *Store) GetTeam(id string) (domain.Team, error) {
	var v domain.Team
	err := s.get("teams", id, &v)
	return v, err
}
func (s *Store) ListTeams() ([]domain.Team, error) {
	var v []domain.Team
	err := s.list("teams", &v)
	sort.Slice(v, func(i, j int) bool { return strings.ToLower(v[i].Name) < strings.ToLower(v[j].Name) })
	return v, err
}

func (s *Store) SavePlayer(player domain.Player) error { return s.save("players", player.ID, player) }
func (s *Store) GetPlayer(id string) (domain.Player, error) {
	var v domain.Player
	err := s.get("players", id, &v)
	return v, err
}
func (s *Store) ListPlayers() ([]domain.Player, error) {
	var v []domain.Player
	err := s.list("players", &v)
	sort.Slice(v, func(i, j int) bool { return v[i].Number < v[j].Number })
	return v, err
}

func (s *Store) SaveGame(game domain.Game) error { return s.save("games", game.ID, game) }
func (s *Store) GetGame(id string) (domain.Game, error) {
	var v domain.Game
	err := s.get("games", id, &v)
	return v, err
}
func (s *Store) ListGames() ([]domain.Game, error) {
	var v []domain.Game
	err := s.list("games", &v)
	sort.Slice(v, func(i, j int) bool { return v[i].ScheduledAt.Before(v[j].ScheduledAt) })
	return v, err
}

func (s *Store) SaveAudit(audit domain.ScoreAudit) error { return s.save("audits", audit.ID, audit) }
func (s *Store) GetAudit(id string) (domain.ScoreAudit, error) {
	var v domain.ScoreAudit
	err := s.get("audits", id, &v)
	return v, err
}
func (s *Store) ListAudits() ([]domain.ScoreAudit, error) {
	var v []domain.ScoreAudit
	err := s.list("audits", &v)
	sort.Slice(v, func(i, j int) bool { return v[i].ChangedAt.Before(v[j].ChangedAt) })
	return v, err
}

func (s *Store) SaveUser(user domain.UserAccount) error { return s.save("users", user.ID, user) }
func (s *Store) GetUser(id string) (domain.UserAccount, error) {
	var v domain.UserAccount
	err := s.get("users", id, &v)
	return v, err
}
func (s *Store) ListUsers() ([]domain.UserAccount, error) {
	var v []domain.UserAccount
	err := s.list("users", &v)
	return v, err
}

func (s *Store) SaveMedia(media domain.MediaAsset) error { return s.save("media", media.ID, media) }
func (s *Store) GetMedia(id string) (domain.MediaAsset, error) {
	var v domain.MediaAsset
	err := s.get("media", id, &v)
	return v, err
}
func (s *Store) ListMedia() ([]domain.MediaAsset, error) {
	var v []domain.MediaAsset
	err := s.list("media", &v)
	sort.Slice(v, func(i, j int) bool { return v[i].CapturedAt.After(v[j].CapturedAt) })
	return v, err
}

func (s *Store) SaveMessage(message domain.GuestMessage) error {
	return s.save("messages", message.ID, message)
}
func (s *Store) GetMessage(id string) (domain.GuestMessage, error) {
	var v domain.GuestMessage
	err := s.get("messages", id, &v)
	return v, err
}
func (s *Store) ListMessages() ([]domain.GuestMessage, error) {
	var v []domain.GuestMessage
	err := s.list("messages", &v)
	sort.Slice(v, func(i, j int) bool { return v[i].CreatedAt.After(v[j].CreatedAt) })
	return v, err
}

func (s *Store) Delete(bucket, id string) error {
	return s.transaction(true, func(tx *bolt.Tx) error { return tx.Bucket([]byte(bucket)).Delete(key(id)) })
}

func (s *Store) save(bucket, id string, value any) error {
	data, err := marshal(value)
	if err != nil {
		return err
	}
	return s.transaction(true, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		return b.Put(key(id), data)
	})
}

func (s *Store) get(bucket, id string, target any) error {
	return s.transaction(false, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		data := b.Get(key(id))
		if data == nil {
			return ErrNotFound
		}
		copyData := append([]byte(nil), data...)
		return unmarshal(copyData, target)
	})
}

func (s *Store) list(bucket string, target any) error {
	return s.transaction(false, func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return ErrNotFound
		}
		rv := reflect.ValueOf(target)
		if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
			return ErrNotFound
		}
		slice := reflect.MakeSlice(rv.Elem().Type(), 0, b.Stats().KeyN)
		err := b.ForEach(func(_, data []byte) error {
			if data == nil {
				return nil
			}
			item := reflect.New(rv.Elem().Type().Elem())
			if err := json.Unmarshal(data, item.Interface()); err != nil {
				return err
			}
			slice = reflect.Append(slice, item.Elem())
			return nil
		})
		if err != nil {
			return err
		}
		rv.Elem().Set(slice)
		return nil
	})
}
