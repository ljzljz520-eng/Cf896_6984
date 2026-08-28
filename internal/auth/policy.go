package auth

import (
	"fmt"
	"strings"

	"townbasketball/internal/domain"
)

type Permission string

const (
	PermissionViewPublic      Permission = "view_public"
	PermissionSubmitRoster    Permission = "submit_roster"
	PermissionReviewRoster    Permission = "review_roster"
	PermissionPublishScore    Permission = "publish_score"
	PermissionModerateMessage Permission = "moderate_message"
)

func Permissions(role domain.Role) []Permission {
	switch role {
	case domain.RoleAdmin:
		return []Permission{PermissionViewPublic, PermissionSubmitRoster, PermissionReviewRoster, PermissionPublishScore, PermissionModerateMessage}
	case domain.RoleCaptain:
		return []Permission{PermissionViewPublic, PermissionSubmitRoster}
	default:
		return []Permission{PermissionViewPublic}
	}
}

func HasPermission(role domain.Role, permission Permission) bool {
	for _, candidate := range Permissions(role) {
		if candidate == permission {
			return true
		}
	}
	return false
}

func ValidateUsername(username string) error {
	value := strings.TrimSpace(username)
	if len(value) < 3 || len(value) > 32 {
		return fmt.Errorf("username must be between 3 and 32 characters")
	}
	for _, r := range value {
		if r == ' ' || r == '\n' || r == '\t' {
			return fmt.Errorf("username cannot contain whitespace")
		}
	}
	return nil
}

func (a *Service) Can(session Session, permission Permission) bool {
	if session.UserID == "" || !HasPermission(session.Role, permission) {
		return false
	}
	return true
}
