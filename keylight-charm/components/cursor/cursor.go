package kl_cursor

import (
	"fmt"
	"keylight-charm/styles"

	"github.com/charmbracelet/lipgloss"
)

func RenderLine(line string, isActive, isEdit bool) string {
	style := lipgloss.NewStyle().PaddingLeft(styles.Padding)
	cursor := ""
	if isActive {
		style = style.UnsetPaddingLeft()
		cursorStyles := lipgloss.NewStyle()
		cursorStyles.Foreground(styles.AccentColor)
		cursor = styles.TextAccentColor.Render(">")
	}
	edit := ""
	if isActive && isEdit {
		edit = "(edit)"
	}
	return style.Render(fmt.Sprintf("%s %s %s", cursor, line, edit))
}
