package hue_lights_details

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	adapter *hue.HueAdapter
	light   *hue_control.Light
}

func InitModel(adapter *hue.HueAdapter, light *hue_control.Light) Model {
	return Model{
		adapter: adapter,
		light:   light,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return m.light.GetName()
}
