package page_help

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type SetKeyMapMsg struct {
	keyMap help.KeyMap
}

func CreateSetKeyMapMsg(keyMap help.KeyMap) tea.Cmd {
	return func() tea.Msg {
		return SetKeyMapMsg{
			keyMap: keyMap,
		}
	}
}

type ResetKeyMapMsg struct{}

func CreateResetKeyMapMsg() tea.Cmd {
	return func() tea.Msg {
		return ResetKeyMapMsg{}
	}
}
