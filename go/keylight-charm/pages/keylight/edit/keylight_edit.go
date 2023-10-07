package keylight_edit

import (
	keylight_model "keylight-charm/pages/keylight/details/model"
	"keylight-control/control"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	keylight *control.Keylight
}

func InitModel(keylight *control.Keylight) Model {
	return Model{keylight}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			cmd = keylight_model.CreateAbortAction()
		}
	}

	return m, cmd
}

func (m Model) View() string {
	return "Edit"
}
