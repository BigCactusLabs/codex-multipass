package app

import (
	"fmt"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/BigCactusLabs/codex-multipass/internal/ui"
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save current auth as a profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fail("Usage: codex-mp save <name>")
		}
		name := args[0]

		paths := config.ResolvePaths()
		profilePath, err := profile.Save(name, paths)
		if err != nil {
			fail(err.Error())
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"save","profile":"%s","path":"%s"}`+"\n", name, profilePath)
		} else {
			ui.Success("Saved profile: %s", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
