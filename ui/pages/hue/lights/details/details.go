package hue_lights_details

import (
	hue_control "hue-control/pubilc"
	"ui/lights/hue"
	hue_lights "ui/pages/hue/lights"

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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = hue_lights.CreateBackToLightHomeAction()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return m.light.GetName()
}
