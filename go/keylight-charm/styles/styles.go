package styles

import "github.com/charmbracelet/lipgloss"

var AccentColor = lipgloss.Color("#1f4a5c")
var WarningColor = lipgloss.Color("#FFA500")
var TextAccentColor = lipgloss.NewStyle().Foreground(AccentColor)
var TextWarningColor = lipgloss.NewStyle().Foreground(WarningColor)

var Padding = 1
