package kl_cursor

import (
	"fmt"

	"github.com/remshams/device-control/tui/styles"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

type KeyMap struct {
	Up   key.Binding
	Down key.Binding
}

var CursorKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "down"),
	),
}

type CursorState struct {
	length int
	index  int
}

func InitCursorState(length int) CursorState {
	return CursorState{
		length: length,
		index:  0,
	}
}

func (cursorState CursorState) Index() int {
	return cursorState.index
}

func (cursorState *CursorState) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, CursorKeyMap.Up):
			cursorState.index--
			cursorState.Normalize()
		case key.Matches(msg, CursorKeyMap.Down):
			cursorState.index++
			cursorState.Normalize()
		}
	}
}

func (cursorState *CursorState) Normalize() {
	if cursorState.index < 0 {
		cursorState.index = cursorState.length - 1
	}
	if cursorState.index >= cursorState.length {
		cursorState.index = 0
	}
}
