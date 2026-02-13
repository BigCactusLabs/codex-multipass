package app

import (
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up profiles directory",
	Run: func(cmd *cobra.Command, args []string) {
		paths := config.ResolvePaths()

		// Create directories with 0700 permissions
		// Note: MkdirAll uses the permission specific for the final directory, 
		// but intermediate directories rely on umask.
		// For security, we explicitly chmod the critical ones.

		if err := os.MkdirAll(paths.ProfilesDir, 0700); err != nil {
			fail("Failed to create profiles directory: %v", err)
		}

		// Enforce 0700 on CODEX_DIR and PROFILES_DIR explicitly
		if err := os.Chmod(paths.CodexDir, 0700); err != nil {
			fail("Failed to set permissions on %s: %v", paths.CodexDir, err)
		}
		if err := os.Chmod(paths.ProfilesDir, 0700); err != nil {
			fail("Failed to set permissions on %s: %v", paths.ProfilesDir, err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"init","profiles_dir":"%s"}`+"\n", paths.ProfilesDir)
		} else {
			fmt.Printf("✓ Initialized profiles directory: %s\n", paths.ProfilesDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
