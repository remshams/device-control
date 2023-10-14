package keylight_model

import (
	"keylight-control/control"

	tea "github.com/charmbracelet/bubbletea"
)

type ViewState string

const (
	Navigate ViewState = "navigate"
	Insert             = "insert"
	isError            = "inError"
)

type UpdateKeylight struct{}

func CreateUpdateKeylight() tea.Cmd {
	return func() tea.Msg {
		return UpdateKeylight{}
	}
}

type CommandStatus string

const (
	NoCommand CommandStatus = "noCommand"
	Success                 = "success"
	Error                   = "error"
)

type CommandResult struct {
	Status CommandStatus
}

func CreateCommandResult(status CommandStatus) tea.Cmd {
	return func() tea.Msg {
		return CommandResult{status}
	}
}

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
