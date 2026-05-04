package app

import (
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/spf13/cobra"
)

var whoCmd = &cobra.Command{
	Use:   "who",
	Short: "Show current auth fingerprint",
	Run: func(cmd *cobra.Command, args []string) {
		paths := config.ResolvePaths()

		fingerprint, err := profile.CurrentAuthFingerprint(paths)
		if err != nil {
			fail("%v", err)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			out := map[string]any{
				"ok":          true,
				"fingerprint": fingerprint,
			}
			if err := writeJSON(os.Stdout, out); err != nil {
				fail("Failed to encode who result: %v", err)
			}
		} else {
			fmt.Println(fingerprint)
		}
	},
}

func init() {
	rootCmd.AddCommand(whoCmd)
}
