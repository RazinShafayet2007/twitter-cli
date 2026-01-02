package cmd

import (
	"fmt"
	"time"

	"github.com/RazinShafayet2007/twitter-cli/internal/config"
	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var blockCmd = &cobra.Command{
	Use:   "block [username]",
	Short: "Block a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetUsername := args[0]

		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in")
		}

		userStore := store.NewUserStore(DB)

		blocker, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		blocked, err := userStore.GetByUsername(targetUsername)
		if err != nil {
			return fmt.Errorf("user @%s not found", targetUsername)
		}

		if blocker.ID == blocked.ID {
			return fmt.Errorf("you cannot block yourself")
		}

		// Insert block
		query := `INSERT OR IGNORE INTO blocks (blocker_id, blocked_id, created_at) VALUES (?, ?, ?)`
		_, err = DB.Exec(query, blocker.ID, blocked.ID, time.Now().Unix())
		if err != nil {
			return fmt.Errorf("failed to block user: %w", err)
		}

		fmt.Printf("Blocked @%s\n", targetUsername)
		return nil
	},
}

var unblockCmd = &cobra.Command{
	Use:   "unblock [username]",
	Short: "Unblock a user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetUsername := args[0]

		username, err := config.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("not logged in")
		}

		userStore := store.NewUserStore(DB)

		blocker, err := userStore.GetByUsername(username)
		if err != nil {
			return err
		}

		blocked, err := userStore.GetByUsername(targetUsername)
		if err != nil {
			return fmt.Errorf("user @%s not found", targetUsername)
		}

		// Remove block
		query := `DELETE FROM blocks WHERE blocker_id = ? AND blocked_id = ?`
		result, err := DB.Exec(query, blocker.ID, blocked.ID)
		if err != nil {
			return fmt.Errorf("failed to unblock user: %w", err)
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			return fmt.Errorf("user @%s was not blocked", targetUsername)
		}

		fmt.Printf("Unblocked @%s\n", targetUsername)
		return nil
	},
}

var blockListCmd = &cobra.Command{
	Use:   "blocked",
	Short: "List blocked users",
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

		query := `
			SELECT u.username 
			FROM blocks b
			JOIN users u ON b.blocked_id = u.id
			WHERE b.blocker_id = ?
			ORDER BY u.username
		`

		rows, err := DB.Query(query, user.ID)
		if err != nil {
			return fmt.Errorf("failed to get blocked users: %w", err)
		}
		defer rows.Close()

		var blockedUsers []string
		for rows.Next() {
			var username string
			if err := rows.Scan(&username); err != nil {
				return err
			}
			blockedUsers = append(blockedUsers, username)
		}

		if len(blockedUsers) == 0 {
			fmt.Println("You haven't blocked anyone.")
			return nil
		}

		fmt.Println("Blocked users:")
		for _, u := range blockedUsers {
			fmt.Printf("  @%s\n", u)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(blockCmd)
	rootCmd.AddCommand(unblockCmd)
	rootCmd.AddCommand(blockListCmd)
}
