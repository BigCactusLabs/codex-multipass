package app

import (
	"fmt"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a profile",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fail("Usage: codex-mp rename <old> <new>")
		}
		oldName := args[0]
		newName := args[1]

		paths := config.ResolvePaths()
		err := profile.Rename(oldName, newName, paths)
		if err != nil {
			fail(err.Error())
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			fmt.Printf(`{"ok":true,"action":"rename","old":"%s","new":"%s"}`+"\n", oldName, newName)
		} else {
			fmt.Printf("→ Renamed: %s → %s\n", oldName, newName)
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
