package hue_lights_home

import (
	"keylight-charm/lights/hue"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	list viewState = "list"
)

type Model struct {
	adapter *hue.HueAdapter
	state   viewState
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		state:   list,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return "Lights"
}
