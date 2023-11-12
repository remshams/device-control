package hue_home

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"
	pages_hue "keylight-charm/pages/hue"
	hue_groups "keylight-charm/pages/hue/groups"
	hue_group_details "keylight-charm/pages/hue/groups/details"
	hue_group_list "keylight-charm/pages/hue/groups/list"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	initial viewState = "init"
	list    viewState = "list"
	details viewState = "details"
)

type initMsg struct{}

type Model struct {
	adapter *hue.HueAdapter
	bridges []hue_control.Bridge
	state   viewState
	list    hue_group_list.Model
	details hue_group_details.Model
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
		m.reloadLights()
		m.bridges = m.adapter.Control.GetBridges()
		m.list = hue_group_list.InitModel(m.adapter, m.bridges)
		m.state = list
	case pages_hue.ReloadBridgesAction:
		m.reloadLights()
	case hue_group_list.GroupSelect:
		m.reloadLights()
		m.details = hue_group_details.InitModel(m.adapter, *m.bridges[0].FindGroup(msg.Group.GetId()))
		m.state = details
	case hue_groups.BackToListAction:
		m.state = list
	default:
		cmd = m.forwardUpdate(msg)
	}
	return m, cmd
}

func (m *Model) forwardUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case list:
		m.list, cmd = m.list.Update(msg)
	case details:
		m.details, cmd = m.details.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case initial:
		return "Loading..."
	case list:
		return m.list.View()
	case details:
		return m.details.View()
	default:
		return ""
	}
}

func (m *Model) reloadLights() {
	m.adapter.Control.LoadOrFindBridges()
	m.bridges = m.adapter.Control.GetBridges()
}

func (m *Model) init() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}
