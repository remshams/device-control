package hue_lights_home

import (
	"keylight-charm/lights/hue"
	hue_lights_list "keylight-charm/pages/hue/lights/list"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	list viewState = "list"
)

type Model struct {
	adapter *hue.HueAdapter
	state   viewState
	list    hue_lights_list.Model
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		state:   list,
		list:    hue_lights_list.InitModel(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case list:
		return m.list.View()
	default:
		return ""
	}
}
