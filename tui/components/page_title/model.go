package page_title

import tea "github.com/charmbracelet/bubbletea"

type SetPageTitleMsg = string

func CreateSetPageTitleMsg(header string) tea.Cmd {
	return func() tea.Msg {
		return SetPageTitleMsg(header)
	}
}
