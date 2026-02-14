package app

import (
	"encoding/json"
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

		fingerprint, err := profile.GetFingerprint(paths.AuthFile)
		if err != nil {
			fail("Not logged in (missing %s)", paths.AuthFile)
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			out := map[string]any{
				"ok":          true,
				"fingerprint": fingerprint,
			}
			json.NewEncoder(os.Stdout).Encode(out)
		} else {
			fmt.Println(fingerprint)
		}
	},
}

func init() {
	rootCmd.AddCommand(whoCmd)
}
