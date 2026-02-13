package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/fs"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch to a saved profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fail("Usage: codex-mp use <name>")
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

		// Atomic Copy Profile -> Auth
		if err := fs.AtomicCopy(profilePath, paths.AuthFile, 0600); err != nil {
			fail("Failed to switch profile: %v", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"use","profile":"%s","auth":"%s"}`+"\n", name, paths.AuthFile)
		} else {
			fmt.Printf("⚡ Switched -> %s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
