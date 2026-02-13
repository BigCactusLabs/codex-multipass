package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/fs"
	"github.com/spf13/cobra"
)

var nameRegex = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

var saveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current auth as a profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fail("Usage: codex-mp save <name>")
		}
		name := args[0]

		if !nameRegex.MatchString(name) {
			fail("Invalid profile name: %s (allowed: A-Z a-z 0-9 . _ -)", name)
		}

		paths := config.ResolvePaths()

		// Read Auth - Verify existence only
		if _, err := os.Stat(paths.AuthFile); os.IsNotExist(err) {
			fail("Missing auth file: %s. Hint: run 'codex login' first.", paths.AuthFile)
		}

		// Save Profile
		profilePath := filepath.Join(paths.ProfilesDir, name+".json")

		// Acquire Lock
		unlock, err := fs.Lock(filepath.Join(paths.CodexDir, ".codex-mp.lock"))
		if err != nil {
			fail("Failed to acquire lock: %v", err)
		}
		defer unlock()

		// Atomic Copy (Raw bytes, just like Bash version)
		if err := fs.AtomicCopy(paths.AuthFile, profilePath, 0600); err != nil {
			fail("Failed to save profile: %v", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"save","profile":"%s","path":"%s"}`+"\n", name, profilePath)
		} else {
			fmt.Printf("✓ Saved profile: %s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
