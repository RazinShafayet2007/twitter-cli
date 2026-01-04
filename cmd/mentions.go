package cmd

import (
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/parser"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var mentionsCmd = &cobra.Command{
	Use:   "mentions",
	Short: "View posts that mention you",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")

		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		mentionStore := store.NewMentionStore(DB)
		posts, err := mentionStore.GetMentions(user.ID, limit)
		if err != nil {
			return err
		}

		if len(posts) == 0 {
			fmt.Println("No one has mentioned you yet.")
			return nil
		}

		fmt.Printf("Posts mentioning @%s:\n\n", username)

		for _, pwa := range posts {
			timeAgo := display.FormatTimeAgo(pwa.Post.CreatedAt)
			fmt.Printf("%s  @%s  %s\n", pwa.Post.ID, pwa.Username, timeAgo)
			fmt.Printf("%s\n", parser.HighlightText(pwa.Post.Text))
			fmt.Println()
		}

		return nil
	},
}

func init() {
	mentionsCmd.Flags().Int("limit", 50, "Number of mentions to show")
	rootCmd.AddCommand(mentionsCmd)
}
