package hue_bridges_home

import (
	"github.com/remshams/device-control/tui/lights/hue"
	pages_hue "github.com/remshams/device-control/tui/pages/hue"
	hue_bridges "github.com/remshams/device-control/tui/pages/hue/bridges"
	hue_bridges_list "github.com/remshams/device-control/tui/pages/hue/bridges/list"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	adapter *hue.HueAdapter
	list    hue_bridges_list.Model
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		list:    hue_bridges_list.InitModel(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.list.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case hue_bridges.BackToBridgesHomeAction:
		cmd = pages_hue.CreateBackToHueHomeAction()
	default:
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
