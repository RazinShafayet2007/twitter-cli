package cmd

import (
	"fmt"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/parser"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var hashtagCmd = &cobra.Command{
	Use:   "hashtag [tag]",
	Short: "View posts with a hashtag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tag := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		// Remove # if user included it
		if len(tag) > 0 && tag[0] == '#' {
			tag = tag[1:]
		}

		hashtagStore := store.NewHashtagStore(DB)
		posts, err := hashtagStore.GetPostsByHashtag(tag, limit)
		if err != nil {
			return err
		}

		if len(posts) == 0 {
			fmt.Printf("No posts found with #%s\n", tag)
			return nil
		}

		fmt.Printf("Posts with #%s:\n\n", tag)

		// Display posts with highlighted text
		for _, pwa := range posts {
			timeAgo := display.FormatTimeAgo(pwa.Post.CreatedAt)
			fmt.Printf("%s  @%s  %s\n", pwa.Post.ID, pwa.Username, timeAgo)
			fmt.Printf("%s\n", parser.HighlightText(pwa.Post.Text))
			fmt.Println()
		}

		return nil
	},
}

var trendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "View trending hashtags",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		days, _ := cmd.Flags().GetInt("days")

		// Calculate "since" timestamp
		since := time.Now().AddDate(0, 0, -days).Unix()

		hashtagStore := store.NewHashtagStore(DB)
		trending, err := hashtagStore.GetTrendingHashtags(limit, since)
		if err != nil {
			return err
		}

		if len(trending) == 0 {
			fmt.Println("No trending hashtags found.")
			return nil
		}

		fmt.Printf("Trending hashtags (last %d days):\n\n", days)

		for i, t := range trending {
			fmt.Printf("%d. #%s (%d posts)\n", i+1, t.Tag, t.Count)
		}

		return nil
	},
}

func init() {
	hashtagCmd.Flags().Int("limit", 50, "Number of posts to show")
	trendingCmd.Flags().Int("limit", 10, "Number of hashtags to show")
	trendingCmd.Flags().Int("days", 7, "Look back this many days")

	rootCmd.AddCommand(hashtagCmd)
	rootCmd.AddCommand(trendingCmd)
}
