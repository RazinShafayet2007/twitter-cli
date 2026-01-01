package validation

import (
	"errors"
	"regexp"
	"strings"
)

const (
	MaxPostLength     = 280
	MaxUsernameLength = 15
	MinUsernameLength = 3
)

// ValidateUsername checks if a username is valid
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)

	if len(username) < MinUsernameLength {
		return errors.New("username must be at least 3 characters")
	}

	if len(username) > MaxUsernameLength {
		return errors.New("username must be at most 15 characters")
	}

	// Only alphanumeric and underscore
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	return nil
}

// ValidatePostText checks if post text is valid
func ValidatePostText(text string) error {
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		return errors.New("post cannot be empty")
	}

	if len(text) > MaxPostLength {
		return errors.New("post cannot exceed 280 characters")
	}

	return nil
}

// SanitizeUsername cleans and lowercases username
func SanitizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// SanitizePostText cleans post text
func SanitizePostText(text string) string {
	return strings.TrimSpace(text)
}
