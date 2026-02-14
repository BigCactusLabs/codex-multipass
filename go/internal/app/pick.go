package app

import (
	"fmt"
	"os"

	"github.com/BigCactusLabs/codex-multipass/internal/config"
	"github.com/BigCactusLabs/codex-multipass/internal/profile"
	"github.com/BigCactusLabs/codex-multipass/internal/ui"
	"github.com/charmbracelet/bubbles/key"
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

		// List profiles
		profiles, err := profile.List(paths)
		if err != nil {
			fail(err.Error())
		}

		if len(profiles) == 0 {
			fail("No profiles found.")
		}

		var options []huh.Option[string]

		for _, p := range profiles {
			label := p.Name
			if p.Active {
				label = fmt.Sprintf("%s (active)", p.Name)
			}
			options = append(options, huh.NewOption(label, p.Name))
		}

		var selectedProfile string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Pick a profile").
					Options(options...).
					Value(&selectedProfile),
			),
		).WithTheme(ui.CustomTheme())

		form.WithKeyMap(func() *huh.KeyMap {
			km := huh.NewDefaultKeyMap()
			km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))
			return km
		}())

		if err := form.Run(); err != nil {
			if err == huh.ErrUserAborted {
				os.Exit(0)
			}
			fail("Error: %v", err)
		}

		if selectedProfile != "" {
			if err := profile.Use(selectedProfile, paths); err != nil {
				fail(err.Error())
			}
			fmt.Printf("⚡ Switched -> %s\n", selectedProfile)
		}
	},
}

func init() {
	rootCmd.AddCommand(pickCmd)
}
