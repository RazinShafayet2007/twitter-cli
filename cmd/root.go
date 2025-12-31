package cmd

import (
	"database/sql"
	"fmt"

	"github.com/RazinShafayet2007/twitter-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	dbPath string
	DB     *sql.DB
)

var rootCmd = &cobra.Command{
	Use:   "twt",
	Short: "A Twitter-like CLI application",
	Long:  `Twitter CLI - A command-line Twitter clone for learning system design`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize database for all commands
		var err error
		DB, err = db.InitDB(dbPath)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		// Close database connection after command
		if DB != nil {
			return DB.Close()
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Default database path
	defaultPath := db.GetDefaultDBPath()

	// Global flags
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", defaultPath, "database file path")
}
