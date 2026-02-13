package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/fs"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fail("Usage: codex-switch rename <old> <new>")
		}
		oldName := args[0]
		newName := args[1]

		if !nameRegex.MatchString(oldName) {
			fail("Invalid profile name: %s", oldName)
		}
		if !nameRegex.MatchString(newName) {
			fail("Invalid profile name: %s", newName)
		}

		paths := config.ResolvePaths()
		oldPath := filepath.Join(paths.ProfilesDir, oldName+".json")
		newPath := filepath.Join(paths.ProfilesDir, newName+".json")

		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			fail("Profile not found: %s", oldName)
		}
		if _, err := os.Stat(newPath); err == nil {
			fail("Profile already exists: %s", newName)
		}

		// Acquire Lock
		unlock, err := fs.Lock(filepath.Join(paths.CodexDir, ".codex-switch.lock"))
		if err != nil {
			fail("Failed to acquire lock: %v", err)
		}
		defer unlock()

		// Rename
		if err := os.Rename(oldPath, newPath); err != nil {
			fail("Failed to rename profile: %v", err)
		}
		
		// Ensure permissions (though rename preserves them usually)
		os.Chmod(newPath, 0600)

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"rename","old":"%s","new":"%s"}`+"\n", oldName, newName)
		} else {
			fmt.Printf("→ Renamed: %s → %s\n", oldName, newName)
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
