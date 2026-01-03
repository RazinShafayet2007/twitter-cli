package cmd

import (
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/RazinShafayet2007/twitter-cli/internal/validation"
	"github.com/spf13/cobra"
)

var userCreateCmd = &cobra.Command{
	Use:   "create [username]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := validation.SanitizeUsername(args[0])

		// Validate username
		if err := validation.ValidateUsername(username); err != nil {
			return err
		}

		userStore := store.NewUserStore(DB)
		user, err := userStore.Create(username)
		if err != nil {
			return err
		}

		fmt.Printf("User @%s created (ID: %s)\n", user.Username, user.ID)
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Login as a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := validation.SanitizeUsername(args[0])

		// Check if user exists
		userStore := store.NewUserStore(DB)
		_, err := userStore.GetByUsername(username)
		if err != nil {
			return fmt.Errorf("user @%s not found", username)
		}

		// Save to config
		if err := config.SetCurrentUser(username); err != nil {
			return fmt.Errorf("failed to save login state: %w", err)
		}

		fmt.Printf("Logged in as @%s\n", username)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentUser, err := config.GetCurrentUser()
		if err != nil {
			return err
		}

		if err := config.ClearCurrentUser(); err != nil {
			return fmt.Errorf("failed to logout: %w", err)
		}

		fmt.Printf("Logged out @%s\n", currentUser)
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current logged-in user",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := config.GetCurrentUser()
		if err != nil {
			return err
		}

		fmt.Printf("@%s\n", username)

		userStore := store.NewUserStore(DB)
		user, err := userStore.GetByUsername(username)
		if err == nil {
			// Show unread messages
			messageStore := store.NewMessageStore(DB)
			unreadMessages, err := messageStore.GetUnreadCount(user.ID)
			if err == nil && unreadMessages > 0 {
				fmt.Printf("💬 %d unread message(s)\n", unreadMessages)
			}

			// Show unread notifications
			notifStore := store.NewNotificationStore(DB)
			unreadNotifs, err := notifStore.GetUnreadCount(user.ID)
			if err == nil && unreadNotifs > 0 {
				fmt.Printf("🔔 %d unread notification(s)\n", unreadNotifs)
			}
		}

		return nil
	},
}

func init() {
	// Create parent 'user' command
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	// Add subcommands
	userCmd.AddCommand(userCreateCmd)

	// Add to root
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
}
