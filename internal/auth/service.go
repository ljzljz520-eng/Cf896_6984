package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"townbasketball/internal/domain"
	"townbasketball/internal/store"
)

type Service struct{ store *store.Store }

type Session struct {
	UserID   string
	Username string
	Role     domain.Role
	TeamID   string
}

func NewService(s *store.Store) *Service { return &Service{store: s} }

func HashPassword(password string) string {
	digest := sha256.Sum256([]byte(password))
	return hex.EncodeToString(digest[:])
}

func (a *Service) Register(user domain.UserAccount, password string) error {
	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(password) == "" {
		return fmt.Errorf("username and password are required")
	}
	if err := domain.ValidateRole(user.Role); err != nil {
		return err
	}
	if user.ID == "" {
		return fmt.Errorf("user id is required")
	}
	user.PasswordHash = HashPassword(password)
	user.Active = true
	return a.store.SaveUser(user)
}

func (a *Service) Login(username, password string) (Session, error) {
	users, err := a.store.ListUsers()
	if err != nil {
		return Session{}, err
	}
	for _, user := range users {
		if user.Username == username && user.Active && user.PasswordHash == HashPassword(password) {
			return Session{UserID: user.ID, Username: user.Username, Role: user.Role, TeamID: user.TeamID}, nil
		}
	}
	return Session{}, fmt.Errorf("invalid credentials")
}

func (a *Service) RequireAdmin(session Session) error {
	if session.Role != domain.RoleAdmin || session.UserID == "" {
		return fmt.Errorf("administrator role required")
	}
	return nil
}

func (a *Service) RequireCaptain(session Session, teamID string) error {
	if session.Role != domain.RoleCaptain || session.TeamID != teamID || session.UserID == "" {
		return fmt.Errorf("team captain role required")
	}
	return nil
}

func (a *Service) SeedDefaults() error {
	users, err := a.store.ListUsers()
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return nil
	}
	defaults := []domain.UserAccount{
		{ID: "admin-001", Username: "admin", Role: domain.RoleAdmin, Active: true, PasswordHash: HashPassword("admin-pass")},
		{ID: "captain-001", Username: "captain-east", Role: domain.RoleCaptain, TeamID: "team-east", Active: true, PasswordHash: HashPassword("captain-pass")},
	}
	for _, user := range defaults {
		if err := a.store.SaveUser(user); err != nil {
			return err
		}
	}
	return nil
}
