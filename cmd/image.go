package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/RazinShafayet2007/twitter-cli/internal/store"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Manage post images",
}

var imageDownloadCmd = &cobra.Command{
	Use:   "download [post_id]",
	Short: "Download images from a post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		postID := args[0]
		outputDir, _ := cmd.Flags().GetString("output")

		mediaStore := store.NewMediaStore(DB)
		mediaList, err := mediaStore.GetByPostID(postID)
		if err != nil {
			return err
		}

		if len(mediaList) == 0 {
			fmt.Println("No images attached to this post.")
			return nil
		}

		// Create output directory
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Copy each image
		for _, m := range mediaList {
			destPath := filepath.Join(outputDir, m.FileName)

			// Copy file
			input, err := os.ReadFile(m.FilePath)
			if err != nil {
				fmt.Printf("Warning: failed to read %s: %v\n", m.FileName, err)
				continue
			}

			if err := os.WriteFile(destPath, input, 0644); err != nil {
				fmt.Printf("Warning: failed to write %s: %v\n", destPath, err)
				continue
			}

			fmt.Printf("Downloaded: %s\n", destPath)
		}

		fmt.Printf("\n✓ Downloaded %d image(s) to %s\n", len(mediaList), outputDir)
		return nil
	},
}

var imageViewCmd = &cobra.Command{
	Use:   "view [post_id]",
	Short: "Open images from a post in default viewer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		postID := args[0]

		mediaStore := store.NewMediaStore(DB)
		mediaList, err := mediaStore.GetByPostID(postID)
		if err != nil {
			return err
		}

		if len(mediaList) == 0 {
			fmt.Println("No images attached to this post.")
			return nil
		}

		// Open each image
		for _, m := range mediaList {
			if err := openFile(m.FilePath); err != nil {
				fmt.Printf("Warning: failed to open %s: %v\n", m.FileName, err)
			} else {
				fmt.Printf("Opened: %s\n", m.FileName)
			}
		}

		return nil
	},
}

// openFile opens a file with the default application
func openFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func init() {
	imageDownloadCmd.Flags().String("output", "./downloads", "Output directory for downloaded images")

	imageCmd.AddCommand(imageDownloadCmd)
	imageCmd.AddCommand(imageViewCmd)

	rootCmd.AddCommand(imageCmd)
}
