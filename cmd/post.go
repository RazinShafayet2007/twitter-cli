package cmd

import (
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/parser"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/RazinShafayet2007/twitter-cli/internal/validation"
	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post [text]",
	Short: "Create a new post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := validation.SanitizePostText(args[0])

		// Validate post text
		if err := validation.ValidatePostText(text); err != nil {
			return err
		}

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

		// Extract and save hashtags
		hashtags := parser.ExtractHashtags(text)
		if len(hashtags) > 0 {
			hashtagStore := store.NewHashtagStore(DB)
			if err := hashtagStore.LinkPostToHashtags(post.ID, hashtags); err != nil {
				fmt.Printf("Warning: failed to save hashtags: %v\n", err)
			}
		}

		// Extract and save mentions
		mentionUsernames := parser.ExtractMentions(text)
		if len(mentionUsernames) > 0 {
			mentionStore := store.NewMentionStore(DB)

			// Get user IDs for mentioned usernames
			mentionedUserIDs, err := mentionStore.GetMentionedUsers(mentionUsernames)
			if err != nil {
				fmt.Printf("Warning: failed to process mentions: %v\n", err)
			} else {
				// Create mention records
				if err := mentionStore.CreateMentions(post.ID, mentionedUserIDs); err != nil {
					fmt.Printf("Warning: failed to save mentions: %v\n", err)
				}

				// Create notifications for mentioned users
				notifStore := store.NewNotificationStore(DB)
				for _, mentionedUserID := range mentionedUserIDs {
					// Don't notify yourself
					if mentionedUserID != user.ID {
						postID := post.ID
						if err := notifStore.Create(mentionedUserID, user.ID, "mention", &postID); err != nil {
							fmt.Printf("Warning: failed to create mention notification: %v\n", err)
						}
					}
				}
			}
		}

		// Show summary
		fmt.Printf("Posted: %s\n", post.ID)
		if len(hashtags) > 0 {
			fmt.Printf("Hashtags: %v\n", hashtags)
		}
		if len(mentionUsernames) > 0 {
			fmt.Printf("Mentions: %v\n", mentionUsernames)
		}

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
		socialStore := store.NewSocialStore(DB)

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

		// Get engagement stats
		likeCount, err := socialStore.GetLikeCount(postID)
		if err != nil {
			return err
		}

		retweetCount, err := postStore.GetRetweetCount(postID)
		if err != nil {
			return err
		}

		// Create PostWithAuthor
		pwa := store.PostWithAuthor{
			Post:     *post,
			Username: user.Username,
		}

		// Display with stats
		fmt.Println(display.FormatPostWithStats(pwa, likeCount, retweetCount))
		return nil
	},
}

var retweetCmd = &cobra.Command{
	Use:   "retweet [post_id]",
	Short: "Retweet a post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		postID := args[0]

		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		// Get original post to find author
		postStore := store.NewPostStore(DB)
		originalPost, err := postStore.GetByID(postID)
		if err != nil {
			return err
		}

		// Create retweet
		retweet, err := postStore.Retweet(user.ID, postID)
		if err != nil {
			return err
		}

		// Create notification for original author
		notifStore := store.NewNotificationStore(DB)
		if err := notifStore.Create(originalPost.AuthorID, user.ID, "retweet", &postID); err != nil {
			fmt.Printf("Warning: failed to create notification: %v\n", err)
		}

		fmt.Printf("Retweeted: %s\n", retweet.ID)
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search posts by text",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		postStore := store.NewPostStore(DB)
		posts, err := postStore.Search(query, 50)
		if err != nil {
			return err
		}

		if len(posts) == 0 {
			fmt.Printf("No posts found matching '%s'\n", query)
			return nil
		}

		fmt.Printf("Found %d post(s) matching '%s':\n\n", len(posts), query)
		output := display.FormatPosts(posts)
		fmt.Println(output)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(deletePostCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(retweetCmd)
	rootCmd.AddCommand(searchCmd)
}
