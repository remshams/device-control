package actions

import (
	"keylight-control/control"

	tea "github.com/charmbracelet/bubbletea"
)

type AbortAction struct{}

func CreateAbortAction() tea.Cmd {
	return func() tea.Msg {
		return AbortAction{}
	}
}

type SaveAction struct {
	Keylight *control.Keylight
}

func CreateSaveAction(keylight *control.Keylight) tea.Cmd {
	return func() tea.Msg {
		return SaveAction{Keylight: keylight}
	}
}
