package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

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
