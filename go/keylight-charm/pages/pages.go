package pages

import tea "github.com/charmbracelet/bubbletea"

type BackToMenuAction struct{}

func CreateBackToMenuAction() tea.Cmd {
	return func() tea.Msg {
		return BackToMenuAction{}
	}
}

type WindowResizeAction struct {
	Width  int
	Height int
}

func CreateWindowResizeAction(width int, height int) tea.Cmd {
	return func() tea.Msg {
		return WindowResizeAction{
			Width:  width,
			Height: height,
		}
	}
}

func IsSystemMsg(msg tea.Msg) bool {
	switch msg.(type) {
	case tea.WindowSizeMsg:
		return true
	}
	return false
}
