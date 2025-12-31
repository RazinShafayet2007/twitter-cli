package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "twt",
	Short: "A Twitter-like CLI application",
	Long:  `Twitter CLI - A command-line Twitter clone for learning system design`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can go here
	// rootCmd.PersistentFlags().StringVar(&dbPath, "db", "~/.twitter-cli/data.db", "database file path")
}
