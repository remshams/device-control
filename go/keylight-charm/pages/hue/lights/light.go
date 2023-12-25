package hue_lights

import tea "github.com/charmbracelet/bubbletea"

type BackToLightHomeAction struct{}

func CreateBackToLightHomeAction() tea.Cmd {
	return func() tea.Msg {
		return BackToLightHomeAction{}
	}
}
