package pages

import tea "github.com/charmbracelet/bubbletea"

type BackToMenuAction struct{}

func CreateBackToMenuAction() tea.Cmd {
	return func() tea.Msg {
		return BackToMenuAction{}
	}
}
