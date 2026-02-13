package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/fs"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fail("Usage: codex-mp delete <name>")
		}
		name := args[0]

		if !nameRegex.MatchString(name) {
			fail("Invalid profile name: %s", name)
		}

		paths := config.ResolvePaths()
		profilePath := filepath.Join(paths.ProfilesDir, name+".json")

		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			fail("Profile not found: %s", name)
		}

		// Acquire Lock
		unlock, err := fs.Lock(filepath.Join(paths.CodexDir, ".codex-mp.lock"))
		if err != nil {
			fail("Failed to acquire lock: %v", err)
		}
		defer unlock()

		// Remove file
		if err := os.Remove(profilePath); err != nil {
			fail("Failed to delete profile: %v", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"delete","profile":"%s"}`+"\n", name)
		} else {
			fmt.Printf("✗ Deleted profile: %s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
