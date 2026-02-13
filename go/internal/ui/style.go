package ui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Color Palette (Cyberpunk/Neon)
var (
	ColorPrimary   = lipgloss.Color("#FF00FF") // Magenta
	ColorSecondary = lipgloss.Color("#00FFFF") // Cyan
	ColorSuccess   = lipgloss.Color("#00FF00") // Neon Green
	ColorError     = lipgloss.Color("#FF0000") // Red
	ColorWarning   = lipgloss.Color("#FFFF00") // Yellow
	ColorText      = lipgloss.Color("#FFFFFF") // White
	ColorSubtext   = lipgloss.Color("#888888") // Grey
)

// Styles
var (
	StyleSuccess = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	StyleError   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	StyleInfo    = lipgloss.NewStyle().Foreground(ColorSecondary)
	StyleWarning = lipgloss.NewStyle().Foreground(ColorWarning)
	StyleTitle   = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Underline(true)
)

// Emojis for Success
var successEmojis = []string{
	"🚀", "✨", "🎉", "🔥", "🌈", "🦄", "🎸", "👾", "🍕", "🍺",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomEmoji returns a random fun emoji
func RandomEmoji() string {
	return successEmojis[rand.Intn(len(successEmojis))]
}

// Success prints a success message with a random emoji
func Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(StyleSuccess.Render(fmt.Sprintf("%s %s", RandomEmoji(), msg)))
}

// Error prints an error message
func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(StyleError.Render(fmt.Sprintf("❌ %s", msg)))
}

// Info prints an info message
func Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(StyleInfo.Render(fmt.Sprintf("ℹ️  %s", msg)))
}

// Warning prints a warning message
func Warning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Println(StyleWarning.Render(fmt.Sprintf("⚠️  %s", msg)))
}

// CustomTheme returns a custom huh theme
func CustomTheme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = lipgloss.NewStyle().Foreground(ColorPrimary)
	t.Focused.Title = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
	t.Focused.Description = lipgloss.NewStyle().Foreground(ColorSubtext)
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	t.Focused.UnselectedOption = lipgloss.NewStyle().Foreground(ColorText)

	t.Blurred.Base = lipgloss.NewStyle().Foreground(ColorSubtext)
	t.Blurred.Title = lipgloss.NewStyle().Foreground(ColorSubtext)
	t.Blurred.Description = lipgloss.NewStyle().Foreground(ColorSubtext)
	t.Blurred.SelectedOption = lipgloss.NewStyle().Foreground(ColorSubtext)
	t.Blurred.UnselectedOption = lipgloss.NewStyle().Foreground(ColorSubtext)

	return t
}
