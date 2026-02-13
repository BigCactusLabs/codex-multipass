package app

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/spf13/cobra"
)

var whoCmd = &cobra.Command{
	Use:   "who",
	Short: "Show current auth fingerprint",
	Run: func(cmd *cobra.Command, args []string) {
		paths := config.ResolvePaths()

		f, err := os.Open(paths.AuthFile)
		if err != nil {
			fail("Not logged in (missing %s)", paths.AuthFile)
		}
		defer f.Close()

		hasher := sha256.New()
		if _, err := io.Copy(hasher, f); err != nil {
			fail("Failed to read auth file: %v", err)
		}

		fingerprint := hex.EncodeToString(hasher.Sum(nil))

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
