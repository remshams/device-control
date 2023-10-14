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

type ErrorAction struct {
	Error string
}

func CreateErrorAction(error string) tea.Cmd {
	return func() tea.Msg {
		return ErrorAction{Error: error}
	}
}
