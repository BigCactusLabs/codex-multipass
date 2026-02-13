package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print tool version",
	Run: func(cmd *cobra.Command, args []string) {
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"version":"%s"}`+"\n", Version)
		} else {
			fmt.Println(Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
