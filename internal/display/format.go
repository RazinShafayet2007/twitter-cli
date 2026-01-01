package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/store"
)

// FormatTimeAgo formats a Unix timestamp as "2m ago", "1h ago", etc.
func FormatTimeAgo(timestamp int64) string {
	now := time.Now().Unix()
	diff := now - timestamp

	if diff < 60 {
		return fmt.Sprintf("%ds ago", diff)
	}
	if diff < 3600 {
		return fmt.Sprintf("%dm ago", diff/60)
	}
	if diff < 86400 {
		return fmt.Sprintf("%dh ago", diff/3600)
	}
	if diff < 604800 {
		return fmt.Sprintf("%dd ago", diff/86400)
	}

	// For older posts, show the date
	t := time.Unix(timestamp, 0)
	return t.Format("Jan 2, 2006")
}

// FormatPost formats a single post for display
func FormatPost(pwa store.PostWithAuthor) string {
	timeAgo := FormatTimeAgo(pwa.Post.CreatedAt)

	var lines []string

	// First line: ID, username, time
	header := fmt.Sprintf("%s  @%s  %s", pwa.Post.ID, pwa.Username, timeAgo)

	// If it's a retweet, show that
	if pwa.Post.IsRetweet {
		lines = append(lines, header)
		lines = append(lines, "↻ Retweeted")
		lines = append(lines, pwa.Post.Text)
	} else {
		lines = append(lines, header)
		lines = append(lines, pwa.Post.Text)
	}

	return strings.Join(lines, "\n")
}

// FormatPosts formats multiple posts for display
func FormatPosts(posts []store.PostWithAuthor) string {
	if len(posts) == 0 {
		return "No posts yet."
	}

	var output []string
	for _, pwa := range posts {
		output = append(output, FormatPost(pwa))
		output = append(output, "") // Empty line between posts
	}

	// Remove trailing empty line
	return strings.Join(output[:len(output)-1], "\n")
}

// FormatPostWithStats formats a post with engagement statistics
func FormatPostWithStats(pwa store.PostWithAuthor, likeCount, retweetCount int) string {
	timeAgo := FormatTimeAgo(pwa.Post.CreatedAt)

	var lines []string

	// Header
	lines = append(lines, fmt.Sprintf("%s  @%s  %s", pwa.Post.ID, pwa.Username, timeAgo))

	// Retweet indicator
	if pwa.Post.IsRetweet {
		lines = append(lines, "↻ Retweeted")
	}

	// Post text
	lines = append(lines, pwa.Post.Text)

	// Engagement stats
	stats := fmt.Sprintf("❤ %d  ↻ %d", likeCount, retweetCount)
	lines = append(lines, stats)

	return strings.Join(lines, "\n")
}
