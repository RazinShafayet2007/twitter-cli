package cmd

import (
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var followCmd = &cobra.Command{
	Use:   "follow [username]",
	Short: "Follow a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetUsername := args[0]

		currentUsername, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)

		currentUser, err := userStore.GetByUsername(currentUsername)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		targetUser, err := userStore.GetByUsername(targetUsername)
		if err != nil {
			return fmt.Errorf("user @%s not found", targetUsername)
		}

		// Follow
		socialStore := store.NewSocialStore(DB)
		err = socialStore.Follow(currentUser.ID, targetUser.ID)
		if err != nil {
			return err
		}

		// Create notification
		notifStore := store.NewNotificationStore(DB)
		if err := notifStore.Create(targetUser.ID, currentUser.ID, "follow", nil); err != nil {
			fmt.Printf("Warning: failed to create notification: %v\n", err)
		}

		fmt.Printf("Now following @%s\n", targetUsername)
		return nil
	},
}

var unfollowCmd = &cobra.Command{
	Use:   "unfollow [username]",
	Short: "Unfollow a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetUsername := args[0]

		// Check if logged in
		currentUsername, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)

		// Get current user ID
		currentUser, err := userStore.GetByUsername(currentUsername)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		// Get target user ID
		targetUser, err := userStore.GetByUsername(targetUsername)
		if err != nil {
			return fmt.Errorf("user @%s not found", targetUsername)
		}

		// Unfollow
		socialStore := store.NewSocialStore(DB)
		err = socialStore.Unfollow(currentUser.ID, targetUser.ID)
		if err != nil {
			return err
		}

		fmt.Printf("Unfollowed @%s\n", targetUsername)
		return nil
	},
}

var followingCmd = &cobra.Command{
	Use:   "following [username]",
	Short: "List users you follow (or another user follows)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetUsername string

		if len(args) == 0 {
			// Show current user's following
			username, err := config.GetCurrentUser()
			if err != nil {
				return fmt.Errorf("not logged in. Run: twt login <username>")
			}
			targetUsername = username
		} else {
			// Show specified user's following
			targetUsername = args[0]
		}

		socialStore := store.NewSocialStore(DB)
		users, err := socialStore.GetFollowingByUsername(targetUsername)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			fmt.Printf("@%s is not following anyone\n", targetUsername)
			return nil
		}

		fmt.Printf("@%s is following:\n", targetUsername)
		for _, user := range users {
			fmt.Printf("  @%s\n", user.Username)
		}

		return nil
	},
}

var followersCmd = &cobra.Command{
	Use:   "followers [username]",
	Short: "List your followers (or another user's followers)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetUsername string

		if len(args) == 0 {
			// Show current user's followers
			username, err := config.GetCurrentUser()
			if err != nil {
				return fmt.Errorf("not logged in. Run: twt login <username>")
			}
			targetUsername = username
		} else {
			// Show specified user's followers
			targetUsername = args[0]
		}

		socialStore := store.NewSocialStore(DB)
		users, err := socialStore.GetFollowersByUsername(targetUsername)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			fmt.Printf("@%s has no followers\n", targetUsername)
			return nil
		}

		fmt.Printf("@%s's followers:\n", targetUsername)
		for _, user := range users {
			fmt.Printf("  @%s\n", user.Username)
		}

		return nil
	},
}

var likeCmd = &cobra.Command{
	Use:   "like [post_id]",
	Short: "Like a post",
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

		// Get post to find author
		postStore := store.NewPostStore(DB)
		post, err := postStore.GetByID(postID)
		if err != nil {
			return err
		}

		// Like the post
		socialStore := store.NewSocialStore(DB)
		if err := socialStore.Like(user.ID, postID); err != nil {
			return err
		}

		// Create notification for post author
		notifStore := store.NewNotificationStore(DB)
		if err := notifStore.Create(post.AuthorID, user.ID, "like", &postID); err != nil {
			// Don't fail if notification fails
			fmt.Printf("Warning: failed to create notification: %v\n", err)
		}

		fmt.Println("Liked")
		return nil
	},
}

var unlikeCmd = &cobra.Command{
	Use:   "unlike [post_id]",
	Short: "Unlike a post",
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

		// Unlike the post
		socialStore := store.NewSocialStore(DB)
		if err := socialStore.Unlike(user.ID, postID); err != nil {
			return err
		}

		fmt.Println("Unliked")
		return nil
	},
}

var likesCmd = &cobra.Command{
	Use:   "likes [post_id]",
	Short: "Show who liked a post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		postID := args[0]

		socialStore := store.NewSocialStore(DB)
		users, err := socialStore.GetLikes(postID)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			fmt.Println("No likes yet")
			return nil
		}

		fmt.Printf("Liked by %d user(s):\n", len(users))
		for _, user := range users {
			fmt.Printf("  @%s\n", user.Username)
		}

		return nil
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats [username]",
	Short: "Show user statistics",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetUsername string

		if len(args) == 0 {
			// Show current user's stats
			username, err := config.GetCurrentUser()
			if err != nil {
				return fmt.Errorf("not logged in. Run: twt login <username>")
			}
			targetUsername = username
		} else {
			targetUsername = args[0]
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(targetUsername)
		if err != nil {
			return fmt.Errorf("user @%s not found", targetUsername)
		}

		// Get post count
		postStore := store.NewPostStore(DB)
		posts, err := postStore.GetByUsername(targetUsername, 10000)
		if err != nil {
			return err
		}

		// Get follow counts
		socialStore := store.NewSocialStore(DB)
		following, followers, err := socialStore.GetFollowCounts(user.ID)
		if err != nil {
			return err
		}

		// Get message stats
		sentQuery := `SELECT COUNT(*) FROM messages WHERE sender_id = ?`
		receivedQuery := `SELECT COUNT(*) FROM messages WHERE receiver_id = ?`

		var sentCount, receivedCount int

		err = DB.QueryRow(sentQuery, user.ID).Scan(&sentCount)
		if err != nil {
			return fmt.Errorf("failed to get sent message count: %w", err)
		}

		err = DB.QueryRow(receivedQuery, user.ID).Scan(&receivedCount)
		if err != nil {
			return fmt.Errorf("failed to get received message count: %w", err)
		}

		// Display stats
		fmt.Printf("@%s\n", targetUsername)
		fmt.Println("─────────────────────────")
		fmt.Printf("Posts:             %d\n", len(posts))
		fmt.Printf("Following:         %d\n", following)
		fmt.Printf("Followers:         %d\n", followers)
		fmt.Printf("Messages sent:     %d\n", sentCount)
		fmt.Printf("Messages received: %d\n", receivedCount)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(followCmd)
	rootCmd.AddCommand(unfollowCmd)
	rootCmd.AddCommand(followingCmd)
	rootCmd.AddCommand(followersCmd)
	rootCmd.AddCommand(likeCmd)
	rootCmd.AddCommand(unlikeCmd)
	rootCmd.AddCommand(likesCmd)
	rootCmd.AddCommand(statsCmd)
}
