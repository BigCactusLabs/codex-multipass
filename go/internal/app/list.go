package app

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved profiles",
	Run: func(cmd *cobra.Command, args []string) {
		paths := config.ResolvePaths()

		profiles, err := profile.List(paths)
		if err != nil {
			fail(err.Error())
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			out := map[string]any{
				"ok":       true,
				"profiles": profiles,
			}
			json.NewEncoder(os.Stdout).Encode(out)
		} else {
			fmt.Println("")
			fmt.Println("  Profiles")
			fmt.Println("  ----------------------------")
			for _, p := range profiles {
				short := ""
				if len(p.Fingerprint) >= 12 {
					short = p.Fingerprint[:12]
				}

				if p.Active {
					fmt.Printf("  ▸ %s  %s  active\n", p.Name, short)
				} else {
					fmt.Printf("    %s  %s\n", p.Name, short)
				}
			}
			fmt.Println("")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
