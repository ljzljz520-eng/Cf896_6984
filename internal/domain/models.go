package domain

import (
	"fmt"
	"strings"
	"time"
)

type Role string

const (
	RoleVisitor Role = "visitor"
	RoleCaptain Role = "captain"
	RoleAdmin   Role = "admin"
)

type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ShortName string    `json:"short_name"`
	Coach     string    `json:"coach"`
	City      string    `json:"city"`
	Bio       string    `json:"bio"`
	Approved  bool      `json:"approved"`
	CreatedAt time.Time `json:"created_at"`
}

type Player struct {
	ID          string `json:"id"`
	TeamID      string `json:"team_id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Position    string `json:"position"`
	Approved    bool   `json:"approved"`
	SubmittedBy string `json:"submitted_by"`
}

type Game struct {
	ID          string    `json:"id"`
	HomeTeamID  string    `json:"home_team_id"`
	AwayTeamID  string    `json:"away_team_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Venue       string    `json:"venue"`
	Status      string    `json:"status"`
	HomeScore   int       `json:"home_score"`
	AwayScore   int       `json:"away_score"`
	Published   bool      `json:"published"`
}

type ScoreAudit struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	OldHome   int       `json:"old_home"`
	OldAway   int       `json:"old_away"`
	NewHome   int       `json:"new_home"`
	NewAway   int       `json:"new_away"`
	ChangedBy string    `json:"changed_by"`
	Reason    string    `json:"reason"`
	ChangedAt time.Time `json:"changed_at"`
}

type UserAccount struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         Role   `json:"role"`
	TeamID       string `json:"team_id"`
	Active       bool   `json:"active"`
}

type MediaAsset struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	Caption    string    `json:"caption"`
	Published  bool      `json:"published"`
	CapturedAt time.Time `json:"captured_at"`
}

type GuestMessage struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
}

type Standing struct {
	TeamID        string `json:"team_id"`
	TeamName      string `json:"team_name"`
	Played        int    `json:"played"`
	Wins          int    `json:"wins"`
	Losses        int    `json:"losses"`
	PointsFor     int    `json:"points_for"`
	PointsAgainst int    `json:"points_against"`
	TablePoints   int    `json:"table_points"`
}

func (t Team) Validate() error {
	if strings.TrimSpace(t.ID) == "" || strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("team id and name are required")
	}
	if len([]rune(t.Name)) > 80 {
		return fmt.Errorf("team name is too long")
	}
	if strings.TrimSpace(t.City) == "" {
		return fmt.Errorf("team city is required")
	}
	return nil
}

func (p Player) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.TeamID) == "" {
		return fmt.Errorf("player id and team are required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("player name is required")
	}
	if p.Number < 0 || p.Number > 99 {
		return fmt.Errorf("player number must be between 0 and 99")
	}
	return nil
}

func (g Game) Validate() error {
	if strings.TrimSpace(g.ID) == "" || strings.TrimSpace(g.HomeTeamID) == "" || strings.TrimSpace(g.AwayTeamID) == "" {
		return fmt.Errorf("game identity is required")
	}
	if g.HomeTeamID == g.AwayTeamID {
		return fmt.Errorf("a team cannot play itself")
	}
	if strings.TrimSpace(g.Venue) == "" {
		return fmt.Errorf("game venue is required")
	}
	if g.HomeScore < 0 || g.AwayScore < 0 {
		return fmt.Errorf("scores cannot be negative")
	}
	return nil
}

func (a ScoreAudit) Validate() error {
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.GameID) == "" {
		return fmt.Errorf("audit identity is required")
	}
	if a.NewHome < 0 || a.NewAway < 0 || a.OldHome < 0 || a.OldAway < 0 {
		return fmt.Errorf("audit scores cannot be negative")
	}
	if strings.TrimSpace(a.ChangedBy) == "" {
		return fmt.Errorf("audit actor is required")
	}
	return nil
}

func (m GuestMessage) Validate() error {
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.Nickname) == "" {
		return fmt.Errorf("message identity is required")
	}
	n := len([]rune(strings.TrimSpace(m.Content)))
	if n < 2 || n > 500 {
		return fmt.Errorf("message content must be between 2 and 500 characters")
	}
	return nil
}

func (u UserAccount) CanManage() bool { return u.Active && u.Role == RoleAdmin }
func (u UserAccount) CanSubmitRoster(teamID string) bool {
	return u.Active && u.Role == RoleCaptain && u.TeamID == teamID
}
