package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show resolved paths",
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		paths := config.ResolvePaths()

		if jsonOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(paths); err != nil {
				fail("Failed to encode paths: %v", err)
			}
			return
		}

		// Check if stdout is a TTY
		stat, _ := os.Stdout.Stat()
		isTerminal := (stat.Mode() & os.ModeCharDevice) != 0

		if !isTerminal {
			// Plain output for scripts
			// echo "CODEX_HOME=${CODEX_HOME:-}"
			// echo "CODEX_DIR=$CODEX_DIR"
			// echo "AUTH=$AUTH"
			// echo "PROFILES_DIR=$PROFILES_DIR"

			valHome := paths.CodexHome
			if valHome == "" {
				// Bash logic: ${CODEX_HOME:-} prints empty if not set
			}
			fmt.Printf("CODEX_HOME=%s\n", valHome)
			fmt.Printf("CODEX_DIR=%s\n", paths.CodexDir)
			fmt.Printf("AUTH=%s\n", paths.AuthFile)
			fmt.Printf("PROFILES_DIR=%s\n", paths.ProfilesDir)
			return
		}

		fmt.Println("  Resolved Paths")
		fmt.Println("  ----------------------------")
		if paths.CodexHome != "" {
			fmt.Printf("  CODEX_HOME     = %s\n", paths.CodexHome)
		} else {
			fmt.Printf("  CODEX_HOME     = (not set)\n")
		}
		fmt.Printf("  CODEX_DIR      = %s\n", paths.CodexDir)
		fmt.Printf("  AUTH           = %s\n", paths.AuthFile)
		fmt.Printf("  PROFILES_DIR   = %s\n", paths.ProfilesDir)
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
