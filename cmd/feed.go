package cmd

import (
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var (
	feedLimit  int
	feedOffset int
)

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "View your personalized feed",
	Long:  `Shows posts from users you follow and your own posts, sorted by time`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if logged in
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		// Get user ID
		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Get feed
		postStore := store.NewPostStore(DB)
		posts, err := postStore.GetFeed(user.ID, feedLimit, feedOffset)
		if err != nil {
			return err
		}

		// Display feed
		if len(posts) == 0 {
			if feedOffset > 0 {
				fmt.Println("No more posts.")
			} else {
				fmt.Println("Your feed is empty. Follow some users and start posting!")
			}
			return nil
		}

		output := display.FormatPosts(posts)
		fmt.Println(output)

		// Show pagination info
		if len(posts) == feedLimit {
			fmt.Printf("\nShowing %d posts. Use --offset %d to see more.\n", feedLimit, feedOffset+feedLimit)
		}

		return nil
	},
}

func init() {
	feedCmd.Flags().IntVar(&feedLimit, "limit", 20, "Number of posts to show")
	feedCmd.Flags().IntVar(&feedOffset, "offset", 0, "Number of posts to skip")

	rootCmd.AddCommand(feedCmd)
}
