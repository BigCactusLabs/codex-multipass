package app

import (
	"fmt"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch to a saved profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fail("Usage: codex-mp use <name>")
		}
		name := args[0]

		paths := config.ResolvePaths()
		err := profile.Use(name, paths)
		if err != nil {
			fail(err.Error())
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"use","profile":"%s","auth":"%s"}`+"\n", name, paths.AuthFile)
		} else {
			fmt.Printf("⚡ Switched -> %s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
