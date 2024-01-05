package header

import tea "github.com/charmbracelet/bubbletea"

type SetHeaderMsg = string

func CreateSetHeaderMsg(header string) tea.Cmd {
	return func() tea.Msg {
		return SetHeaderMsg(header)
	}
}
