package app

import (
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up profiles directory",
	Run: func(cmd *cobra.Command, args []string) {
		paths := config.ResolvePaths()

		if err := profile.EnsureInitialized(paths); err != nil {
			fail("%v", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			if err := writeJSON(os.Stdout, map[string]any{
				"ok":           true,
				"action":       "init",
				"profiles_dir": paths.ProfilesDir,
			}); err != nil {
				fail("Failed to encode init result: %v", err)
			}
		} else {
			fmt.Printf("✓ Initialized profiles directory: %s\n", paths.ProfilesDir)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
