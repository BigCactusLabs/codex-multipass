package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var pickCmd = &cobra.Command{
	Use:     "pick",
	Aliases: []string{"ui"},
	Short:   "Interactive profile selector",
	Run: func(cmd *cobra.Command, args []string) {
		if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
			fail("pick/ui does not support --json output")
			return // unreachable due to fail/exit but good practice
		}
		paths := config.ResolvePaths()

		// list profiles
		entries, err := os.ReadDir(paths.ProfilesDir)
		if err != nil && !os.IsNotExist(err) {
			fail("Failed to list profiles: %v", err)
		}

		var options []huh.Option[string]
		
		// Get Active Fingerprint to mark current
		activeFp, _ := getFingerprint(paths.AuthFile)

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), ".json")
			
			// Check if active
			label := name
			fullPath := filepath.Join(paths.ProfilesDir, entry.Name())
			if fp, _ := getFingerprint(fullPath); fp == activeFp && activeFp != "" {
				label = fmt.Sprintf("%s (active)", name)
			}

			options = append(options, huh.NewOption(label, name))
		}

		if len(options) == 0 {
			fail("No profiles found.")
		}

		var selectedProfile string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Pick a profile").
					Options(options...).
					Value(&selectedProfile),
			),
		)

		if err := form.Run(); err != nil {
			if err == huh.ErrUserAborted {
				os.Exit(0)
			}
			fail("Error: %v", err)
		}

		if selectedProfile != "" {
			// Reuse the Use command logic by calling it directly or via subcommand
			// For simplicity and decoupling, we'll re-run the core logic or exec.
			// Ideally refactor 'use' logic to a shared function.
			// For now, let's just invoke the use command logic via a new args call? 
			// No, better to extract the logic. But we are in 'app' package.
			// Let's just manually run the update since we have the name.
			
			// Actually, Cobra commands are public. We can just call Run.
			// But arguments are passed via args slice.
			useCmd.Run(useCmd, []string{selectedProfile})
		}
	},
}

func init() {
	rootCmd.AddCommand(pickCmd)
}
