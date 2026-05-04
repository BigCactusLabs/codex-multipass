package app

import (
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
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

		paths := config.ResolvePaths()
		err := profile.Delete(name, paths)
		if err != nil {
			fail(err.Error())
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			if err := writeJSON(os.Stdout, map[string]any{
				"ok":      true,
				"action":  "delete",
				"profile": name,
			}); err != nil {
				fail("Failed to encode delete result: %v", err)
			}
		} else {
			fmt.Printf("✗ Deleted profile: %s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
