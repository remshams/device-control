package hue_home

import (
	"fmt"
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"

	tea "github.com/charmbracelet/bubbletea"
)

type initMsg struct {
	Bridges []hue_control.Bridge
}

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

func (m Model) Init() tea.Cmd {
	return m.init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initMsg:
		m.bridges = msg.Bridges
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("HueLights: %d", len(m.bridges))
}

func (m *Model) init() tea.Cmd {
	return func() tea.Msg {
		m.adapter.Control.LoadOrFindBridges()
		return initMsg{
			Bridges: m.adapter.Control.GetBridges(),
		}
	}
}
