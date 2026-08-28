package domain

import (
	"fmt"
	"strings"
	"time"
)

type Announcement struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Category  string    `json:"category"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
}

func (a Announcement) Validate() error {
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.Title) == "" {
		return fmt.Errorf("announcement id and title are required")
	}
	if strings.TrimSpace(a.Body) == "" {
		return fmt.Errorf("announcement body is required")
	}
	if len([]rune(a.Title)) > 120 || len([]rune(a.Body)) > 3000 {
		return fmt.Errorf("announcement is too long")
	}
	if err := EnsureOneOf(a.Category, "schedule", "score", "notice"); err != nil {
		return err
	}
	return nil
}

func (a Announcement) Summary() string {
	body := strings.TrimSpace(a.Body)
	runes := []rune(body)
	if len(runes) > 100 {
		body = string(runes[:100]) + "..."
	}
	return body
}
