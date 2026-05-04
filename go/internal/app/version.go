package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print tool version",
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			if err := writeJSON(os.Stdout, map[string]any{
				"ok":      true,
				"version": Version,
			}); err != nil {
				fail("Failed to encode version result: %v", err)
			}
		} else {
			fmt.Println(Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
