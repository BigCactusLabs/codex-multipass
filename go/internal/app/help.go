package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")) // BigCactus Green

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")). // Light Green
			MarginTop(1).
			MarginBottom(0)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")) // Cyan for commands

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")) // Gray for description

	flagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")) // Gray for flags

	exampleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true)
)

func helpFunc(cmd *cobra.Command, args []string) {
	var b strings.Builder

	// Title / Description
	if cmd.Short != "" {
		b.WriteString(titleStyle.Render(cmd.Short) + "\n")
	}
	if cmd.Long != "" {
		if cmd.Short != "" {
			b.WriteString("\n")
		}
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(cmd.Long) + "\n")
	}

	// Usage
	b.WriteString(headerStyle.Render("USAGE") + "\n")
	if cmd.Runnable() {
		b.WriteString(fmt.Sprintf("  %s\n", cmd.UseLine()))
	}
	if cmd.HasAvailableSubCommands() {
		b.WriteString(fmt.Sprintf("  %s [command]\n", cmd.CommandPath()))
	}

	// Subcommands
	if cmd.HasAvailableSubCommands() {
		b.WriteString(headerStyle.Render("COMMANDS") + "\n")

		// Find longest command name for padding
		maxLen := 0
		for _, c := range cmd.Commands() {
			if !c.IsAvailableCommand() {
				continue
			}
			if len(c.Name()) > maxLen {
				maxLen = len(c.Name())
			}
		}

		for _, c := range cmd.Commands() {
			if !c.IsAvailableCommand() {
				continue
			}

			pad := strings.Repeat(" ", maxLen-len(c.Name())+2)
			cmdName := commandStyle.Render(c.Name())
			desc := descStyle.Render(c.Short)

			b.WriteString(fmt.Sprintf("  %s%s%s\n", cmdName, pad, desc))
		}
	}

	// Flags
	if cmd.HasAvailableLocalFlags() {
		b.WriteString(headerStyle.Render("FLAGS") + "\n")
		b.WriteString(flagStyle.Render(cmd.LocalFlags().FlagUsages()))
	}

	// Global Flags
	if cmd.HasAvailableInheritedFlags() {
		b.WriteString(headerStyle.Render("GLOBAL FLAGS") + "\n")
		b.WriteString(flagStyle.Render(cmd.InheritedFlags().FlagUsages()))
	}

	// Examples
	if cmd.Example != "" {
		b.WriteString(headerStyle.Render("EXAMPLES") + "\n")
		b.WriteString(exampleStyle.Render(cmd.Example) + "\n")
	}

	// Footer
	b.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Use \"codex-mp [command] --help\" for more information about a command.") + "\n")

	fmt.Println(b.String())
}
