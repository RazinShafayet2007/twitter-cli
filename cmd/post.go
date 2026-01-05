package cmd

import (
	"fmt"
	"os"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/media"
	"github.com/RazinShafayet2007/twitter-cli/internal/models"
	"github.com/RazinShafayet2007/twitter-cli/internal/parser"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/RazinShafayet2007/twitter-cli/internal/validation"
	"github.com/spf13/cobra"
)

var (
	postImages []string
)

var postCmd = &cobra.Command{
	Use:   "post [text]",
	Short: "Create a new post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := validation.SanitizePostText(args[0])

		if err := validation.ValidatePostText(text); err != nil {
			return err
		}

		// Validate images
		if len(postImages) > media.MaxImagesPerPost {
			return fmt.Errorf("too many images (max %d)", media.MaxImagesPerPost)
		}

		for _, imgPath := range postImages {
			if err := media.ValidateImage(imgPath); err != nil {
				return fmt.Errorf("invalid image %s: %w", imgPath, err)
			}
		}

		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

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

		// Process images
		if len(postImages) > 0 {
			mediaStore := store.NewMediaStore(DB)

			for i, imgPath := range postImages {
				// Copy image to media directory
				destPath, fileName, err := media.CopyImageToMedia(imgPath, post.ID, i)
				if err != nil {
					fmt.Printf("Warning: failed to copy image %s: %v\n", imgPath, err)
					continue
				}

				// Get image info
				width, height, _ := media.GetImageDimensions(imgPath)
				fileType, _ := media.GetFileType(imgPath)
				fileInfo, _ := os.Stat(imgPath)

				// Create media record
				m := &models.Media{
					PostID:   post.ID,
					FilePath: destPath,
					FileName: fileName,
					FileType: fileType,
					FileSize: fileInfo.Size(),
					Width:    &width,
					Height:   &height,
					Position: i,
				}

				if err := mediaStore.Create(m); err != nil {
					fmt.Printf("Warning: failed to save media record: %v\n", err)
				}
			}
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

			mentionedUserIDs, err := mentionStore.GetMentionedUsers(mentionUsernames)
			if err != nil {
				fmt.Printf("Warning: failed to process mentions: %v\n", err)
			} else {
				if err := mentionStore.CreateMentions(post.ID, mentionedUserIDs); err != nil {
					fmt.Printf("Warning: failed to save mentions: %v\n", err)
				}

				// Create notifications
				notifStore := store.NewNotificationStore(DB)
				for _, mentionedUserID := range mentionedUserIDs {
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
		if len(postImages) > 0 {
			fmt.Printf("📷 %d image(s) attached\n", len(postImages))
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
		mediaStore := store.NewMediaStore(DB)
		for _, pwa := range posts {
			mediaCount, _ := mediaStore.GetMediaCount(pwa.Post.ID)
			fmt.Println(display.FormatPostWithMedia(pwa, mediaCount))
			fmt.Println()
		}

		return nil
	},
}

var deletePostCmd = &cobra.Command{
	Use:   "delete [post_id]",
	Short: "Delete your own post",
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

		// Get media before deleting post
		mediaStore := store.NewMediaStore(DB)
		mediaList, err := mediaStore.GetByPostID(postID)
		if err != nil {
			return err
		}

		// Delete post
		postStore := store.NewPostStore(DB)
		if err := postStore.Delete(postID, user.ID); err != nil {
			return err
		}

		// Delete media files from disk
		for _, m := range mediaList {
			if err := media.DeleteMediaFile(m.FilePath); err != nil {
				fmt.Printf("Warning: failed to delete media file: %v\n", err)
			}
		}

		fmt.Println("Post deleted")
		if len(mediaList) > 0 {
			fmt.Printf("Deleted %d image(s)\n", len(mediaList))
		}

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
		mediaStore := store.NewMediaStore(DB)

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

		// Get media
		mediaList, err := mediaStore.GetByPostID(postID)
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

		// Show media info
		if len(mediaList) > 0 {
			fmt.Printf("\n📷 %d image(s) attached:\n", len(mediaList))
			for i, m := range mediaList {
				sizeKB := m.FileSize / 1024
				fmt.Printf("  %d. %s (%d KB", i+1, m.FileName, sizeKB)
				if m.Width != nil && m.Height != nil {
					fmt.Printf(", %dx%d", *m.Width, *m.Height)
				}
				fmt.Printf(")\n")
			}
			fmt.Printf("\nDownload: twt image download %s\n", postID)
		}

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
	// Add image flag
	postCmd.Flags().StringArrayVar(&postImages, "image", []string{}, "Attach image(s) to post (can be used multiple times)")
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(deletePostCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(retweetCmd)
	rootCmd.AddCommand(searchCmd)
}
