package keylight_model

import tea "github.com/charmbracelet/bubbletea"

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
