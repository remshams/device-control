package page_settings

import tea "github.com/charmbracelet/bubbletea"

type BackToSettingsHomeMsg struct{}

func CreateBackToSettingsHomeAction() tea.Cmd {
	return func() tea.Msg {
		return BackToSettingsHomeMsg{}
	}
}
