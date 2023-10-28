package pages_keylight

import tea "github.com/charmbracelet/bubbletea"

type BackAction struct{}

func CreateBackAction() tea.Cmd {
	return func() tea.Msg {
		return BackAction{}
	}
}

type ReloadKeylights struct{}

func CreateReloadKeylights() tea.Cmd {
	return func() tea.Msg {
		return ReloadKeylights{}
	}
}
