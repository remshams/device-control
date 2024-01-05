package styles

import "github.com/charmbracelet/lipgloss"

var AccentColor = lipgloss.Color("#1f4a5c")
var WarningColor = lipgloss.Color("#FFA500")
var ErrorColor = lipgloss.Color("#FF0000")
var InfoColor = lipgloss.Color("#387DA4")
var SuccessColor = lipgloss.Color("#00FF00")
var TextAccentColor = lipgloss.NewStyle().Foreground(AccentColor)
var TextWarningColor = lipgloss.NewStyle().Foreground(WarningColor)
var TextErrorColor = lipgloss.NewStyle().Foreground(ErrorColor)
var TextInfoColor = lipgloss.NewStyle().Foreground(InfoColor)
var TextSuccessColor = lipgloss.NewStyle().Foreground(SuccessColor)
var SelectedColor = lipgloss.Color("212")
var DisabledColor = lipgloss.Color("238")
var HeadlineBackgroundColor = lipgloss.Color("#5f00ff")

var Padding = 1
var ListStyles = lipgloss.NewStyle().Margin(1, 2)
