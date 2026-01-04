package parser

import (
	"regexp"
	"strings"
)

// ExtractHashtags extracts all hashtags from text
func ExtractHashtags(text string) []string {
	// Regex: # followed by alphanumeric and underscore
	re := regexp.MustCompile(`#(\w+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	var hashtags []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			tag := strings.ToLower(match[1]) // Normalize to lowercase
			if !seen[tag] {
				hashtags = append(hashtags, tag)
				seen[tag] = true
			}
		}
	}

	return hashtags
}

// ExtractMentions extracts all @mentions from text
func ExtractMentions(text string) []string {
	// Regex: @ followed by alphanumeric and underscore
	re := regexp.MustCompile(`@(\w+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	var mentions []string
	seen := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			username := strings.ToLower(match[1]) // Normalize to lowercase
			if !seen[username] {
				mentions = append(mentions, username)
				seen[username] = true
			}
		}
	}

	return mentions
}

// HighlightHashtags adds color/formatting to hashtags
func HighlightHashtags(text string) string {
	re := regexp.MustCompile(`(#\w+)`)
	return re.ReplaceAllString(text, "\033[36m$1\033[0m") // Cyan color
}

// HighlightMentions adds color/formatting to mentions
func HighlightMentions(text string) string {
	re := regexp.MustCompile(`(@\w+)`)
	return re.ReplaceAllString(text, "\033[33m$1\033[0m") // Yellow color
}

// HighlightText highlights both hashtags and mentions
func HighlightText(text string) string {
	text = HighlightHashtags(text)
	text = HighlightMentions(text)
	return text
}
