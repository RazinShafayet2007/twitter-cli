package cmd

import (
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/display"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var (
	notifUnreadOnly bool
	notifLimit      int
)

var notificationsCmd = &cobra.Command{
	Use:     "notifications",
	Aliases: []string{"notifs"},
	Short:   "View notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in. Run: twt login <username>")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		// Get notifications
		notifStore := store.NewNotificationStore(DB)
		notifications, err := notifStore.GetNotifications(user.ID, notifUnreadOnly, notifLimit)
		if err != nil {
			return err
		}

		if len(notifications) == 0 {
			if notifUnreadOnly {
				fmt.Println("No unread notifications.")
			} else {
				fmt.Println("No notifications yet.")
			}
			return nil
		}

		// Display notifications
		if notifUnreadOnly {
			fmt.Println("Unread notifications:")
		} else {
			fmt.Println("Notifications:")
		}
		fmt.Println()

		for _, n := range notifications {
			timeAgo := display.FormatTimeAgo(n.Notification.CreatedAt)
			unreadIndicator := ""
			if !n.Notification.Read {
				unreadIndicator = " 🔴"
			}

			// Format based on type
			var message string
			switch n.Notification.Type {
			case "like":
				if n.TargetText != nil {
					truncated := truncateText(*n.TargetText, 30)
					message = fmt.Sprintf("@%s liked your post: \"%s\"", n.ActorName, truncated)
				} else {
					message = fmt.Sprintf("@%s liked your post", n.ActorName)
				}
			case "retweet":
				if n.TargetText != nil {
					truncated := truncateText(*n.TargetText, 30)
					message = fmt.Sprintf("@%s retweeted your post: \"%s\"", n.ActorName, truncated)
				} else {
					message = fmt.Sprintf("@%s retweeted your post", n.ActorName)
				}
			case "follow":
				message = fmt.Sprintf("@%s followed you", n.ActorName)
			case "message":
				if n.TargetText != nil {
					truncated := truncateText(*n.TargetText, 30)
					message = fmt.Sprintf("@%s sent you a message: \"%s\"", n.ActorName, truncated)
				} else {
					message = fmt.Sprintf("@%s sent you a message", n.ActorName)
				}
			default:
				message = fmt.Sprintf("@%s performed an action", n.ActorName)
			}

			fmt.Printf("%s (%s)%s\n", message, timeAgo, unreadIndicator)
		}

		fmt.Println()
		fmt.Printf("Showing %d notification(s)\n", len(notifications))

		return nil
	},
}

var notificationsReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Mark all notifications as read",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		notifStore := store.NewNotificationStore(DB)
		if err := notifStore.MarkAsRead(user.ID); err != nil {
			return err
		}

		fmt.Println("All notifications marked as read")
		return nil
	},
}

var notificationsCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Show unread notification count",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		notifStore := store.NewNotificationStore(DB)
		count, err := notifStore.GetUnreadCount(user.ID)
		if err != nil {
			return err
		}

		switch count {
		case 0:
			fmt.Println("No unread notifications.")
		case 1:
			fmt.Println("You have 1 unread notification.")
		default:
			fmt.Printf("You have %d unread notifications.\n", count)
		}

		return nil
	},
}

var notificationsClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all read notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in")
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		notifStore := store.NewNotificationStore(DB)
		if err := notifStore.DeleteAllRead(user.ID); err != nil {
			return err
		}

		fmt.Println("Cleared all read notifications")
		return nil
	},
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

func init() {
	// Flags
	notificationsCmd.Flags().BoolVar(&notifUnreadOnly, "unread", false, "Show only unread notifications")
	notificationsCmd.Flags().IntVar(&notifLimit, "limit", 20, "Number of notifications to show")

	// Subcommands
	notificationsCmd.AddCommand(notificationsReadCmd)
	notificationsCmd.AddCommand(notificationsCountCmd)
	notificationsCmd.AddCommand(notificationsClearCmd)

	rootCmd.AddCommand(notificationsCmd)
}
