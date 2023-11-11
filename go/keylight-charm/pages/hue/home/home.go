package hue_home

import (
	"fmt"
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"
	hue_group_list "keylight-charm/pages/hue/groups/list"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	initial viewState = "init"
	list    viewState = "list"
)

type initMsg struct {
	Bridges []hue_control.Bridge
}

type Model struct {
	adapter *hue.HueAdapter
	bridges []hue_control.Bridge
	state   viewState
	list    hue_group_list.Model
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		bridges: []hue_control.Bridge{},
		state:   initial,
	}
}

func (m Model) Init() tea.Cmd {
	return m.init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		m.bridges = msg.Bridges
		m.list = hue_group_list.InitModel(m.adapter, msg.Bridges)
		m.state = list
	case hue_group_list.GroupSelect:
		fmt.Println(msg.Group.GetId())
	default:
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case initial:
		return "Loading..."
	case list:
		return m.list.View()
	default:
		return ""
	}
}

func (m *Model) init() tea.Cmd {
	return func() tea.Msg {
		m.adapter.Control.LoadOrFindBridges()
		return initMsg{
			Bridges: m.adapter.Control.GetBridges(),
		}
	}
}
