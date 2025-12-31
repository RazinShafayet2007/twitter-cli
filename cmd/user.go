package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userCreateCmd = &cobra.Command{
	Use:   "create [username]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		fmt.Printf("TODO: Create user %s\n", username)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Login as a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		fmt.Printf("TODO: Login as %s\n", username)
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current logged-in user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: Show current user")
	},
}

func init() {
	// Create a parent 'user' command
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	// Add subcommands
	userCmd.AddCommand(userCreateCmd)

	// Add to root
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(whoamiCmd)
}
