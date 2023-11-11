package hue_list

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	adapter *hue.HueAdapter
	bridges []hue_control.Bridge
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		bridges: []hue_control.Bridge{},
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "List"
}
