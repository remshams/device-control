package keylight_model

import tea "github.com/charmbracelet/bubbletea"

type ViewState string

const (
	Navigate ViewState = "navigate"
	Insert             = "insert"
	InError            = "inError"
)

type UpdateKeylight struct{}

func CreateUpdateKeylight() tea.Cmd {
	return func() tea.Msg {
		return UpdateKeylight{}
	}
}
