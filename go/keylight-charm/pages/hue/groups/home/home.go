package hue_groups_home

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
	adapter       *hue.HueAdapter
	state         viewState
	list          hue_group_list.Model
	details       hue_group_details.Model
	selectedGroup *hue_control.Group
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
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
		m.list = hue_group_list.InitModel(m.adapter, m.adapter.Control.GetBridges())
		m.state = list
	case hue_groups.BackToGroupHomeAction:
		cmd = pages_hue.CreateBackToHueHomeAction()
	case pages_hue.ReloadBridgesAction:
		m.reloadLights()
		// TODO Add case to update list view
		cmd = hue_groups.CreateGroupReloadedAction(*m.adapter.Control.GetBridges()[0].GetGroupById(m.selectedGroup.GetId()))
	case hue_group_list.GroupSelect:
		m.selectedGroup = &msg.Group
		m.details = hue_group_details.InitModel(m.adapter, msg.Group)
		m.state = details
	case hue_groups.BackToListAction:
		m.selectedGroup = nil
		m.list = hue_group_list.InitModel(m.adapter, m.adapter.Control.GetBridges())
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
	m.adapter.Control.LoadBridges()
}

func (m *Model) init() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}
