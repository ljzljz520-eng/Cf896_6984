package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{2,63}$`)

func ValidateID(value string) error {
	if !idPattern.MatchString(value) {
		return fmt.Errorf("invalid id %q", value)
	}
	return nil
}

func NormalizeName(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func ValidateRole(role Role) error {
	switch role {
	case RoleVisitor, RoleCaptain, RoleAdmin:
		return nil
	default:
		return fmt.Errorf("unknown role %q", role)
	}
}

func EnsureOneOf(value string, allowed ...string) error {
	for _, candidate := range allowed {
		if value == candidate {
			return nil
		}
	}
	return fmt.Errorf("value %q is not allowed", value)
}

func ValidateScore(home, away int) error {
	if home < 0 || away < 0 {
		return fmt.Errorf("score must be non-negative")
	}
	if home > 300 || away > 300 {
		return fmt.Errorf("score exceeds game limit")
	}
	return nil
}
