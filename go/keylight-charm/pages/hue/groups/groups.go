package hue_groups

import tea "github.com/charmbracelet/bubbletea"

type BackToListAction struct{}

func CreateBackToListAction() tea.Cmd {
	return func() tea.Msg {
		return BackToListAction{}
	}
}
