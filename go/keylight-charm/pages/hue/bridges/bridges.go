package hue_bridges

import tea "github.com/charmbracelet/bubbletea"

type BackToBridgesHomeAction struct{}

func CreateBackToBridgesHomeAction() tea.Cmd {
	return func() tea.Msg {
		return BackToBridgesHomeAction{}
	}
}
