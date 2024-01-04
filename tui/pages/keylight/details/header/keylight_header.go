package keylight_header

import (
	"fmt"
	"github.com/remshams/device-control/keylight-control/control"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	keylight *control.Keylight
}

func InitModel(keylight *control.Keylight) Model {
	return Model{keylight: keylight}
}

func (m Model) Update() (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("Update %s", m.keylight.Metadata.Name)
}
