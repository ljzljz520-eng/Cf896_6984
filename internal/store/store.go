package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("record not found")

var buckets = [][]byte{
	[]byte("teams"), []byte("players"), []byte("games"), []byte("audits"),
	[]byte("users"), []byte("media"), []byte("messages"), []byte("announcements"), []byte("settings"),
}

type Store struct {
	db   *bolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("database path is required")
	}
	db, err := bolt.Open(filepath.Clean(path), 0600, &bolt.Options{Timeout: 2 * time.Second, NoSync: false})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db, path: path}
	if err := s.initBuckets(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range buckets {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", name, err)
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Reopen() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db != nil {
		return fmt.Errorf("store is already open")
	}
	db, err := bolt.Open(filepath.Clean(s.path), 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return fmt.Errorf("reopen database: %w", err)
	}
	s.db = db
	return nil
}

func (s *Store) Path() string { return s.path }

func (s *Store) transaction(write bool, fn func(*bolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store is closed")
	}
	if write {
		return s.db.Update(fn)
	}
	return s.db.View(fn)
}
