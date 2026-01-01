package main

import (
	"fmt"
	"os"

	"github.com/RazinShafayet2007/twitter-cli/cmd"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Set version info
	cmd.SetVersionInfo(Version, BuildTime)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
