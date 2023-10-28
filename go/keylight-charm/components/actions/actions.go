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

type ReloadKeylights struct{}

func CreateReloadKeylights() tea.Cmd {
	return func() tea.Msg {
		return ReloadKeylights{}
	}
}
