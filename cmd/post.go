package cmd

import (
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post [text]",
	Short: "Create a new post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]

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

		// Create post
		postStore := store.NewPostStore(DB)
		post, err := postStore.Create(user.ID, text)
		if err != nil {
			return err
		}

		fmt.Printf("Posted: %s\n", post.ID)
		return nil
	},
}

var profileCmd = &cobra.Command{
	Use:   "profile [username]",
	Short: "View a user's posts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		// Check if user exists
		userStore := store.NewUserStore(DB)
		_, err := userStore.GetByUsername(username)
		if err != nil {
			return fmt.Errorf("user @%s not found", username)
		}

		// Get posts
		postStore := store.NewPostStore(DB)
		posts, err := postStore.GetByUsername(username, 50) // Limit to 50 posts
		if err != nil {
			return err
		}

		// Display posts
		output := display.FormatPosts(posts)
		fmt.Println(output)

		return nil
	},
}

var deletePostCmd = &cobra.Command{
	Use:   "delete [post_id]",
	Short: "Delete your own post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		postID := args[0]

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

		// Delete post
		postStore := store.NewPostStore(DB)
		if err := postStore.Delete(postID, user.ID); err != nil {
			return err
		}

		fmt.Println("Post deleted")
		return nil
	},
}

var showCmd = &cobra.Command{
	Use:   "show [post_id]",
	Short: "Show a single post with details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		postID := args[0]

		postStore := store.NewPostStore(DB)
		userStore := store.NewUserStore(DB)

		// Get post
		post, err := postStore.GetByID(postID)
		if err != nil {
			return err
		}

		// Get author
		user, err := userStore.GetByID(post.AuthorID)
		if err != nil {
			return err
		}

		// Create PostWithAuthor
		pwa := store.PostWithAuthor{
			Post:     *post,
			Username: user.Username,
		}

		// Display
		fmt.Println(display.FormatPost(pwa))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(deletePostCmd)
	rootCmd.AddCommand(showCmd)
}
