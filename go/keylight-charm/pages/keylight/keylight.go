package pages_keylight

import tea "github.com/charmbracelet/bubbletea"

type BackToListAction struct{}

func CreateBackToListAction() tea.Cmd {
	return func() tea.Msg {
		return BackToListAction{}
	}
}

type ReloadKeylights struct{}

func CreateReloadKeylights() tea.Cmd {
	return func() tea.Msg {
		return ReloadKeylights{}
	}
}
