package hue_groups_home

import (
	hue_control "github.com/remshams/device-control/hue-control/pubilc"
	"github.com/remshams/device-control/tui/lights/hue"
	pages_hue "github.com/remshams/device-control/tui/pages/hue"
	hue_groups "github.com/remshams/device-control/tui/pages/hue/groups"
	hue_group_details "github.com/remshams/device-control/tui/pages/hue/groups/details"
	hue_group_list "github.com/remshams/device-control/tui/pages/hue/groups/list"

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
		list:    hue_group_list.InitModel(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		m.list = hue_group_list.InitModel(m.adapter)
		m.state = list
	case hue_groups.BackToGroupHomeAction:
		cmd = pages_hue.CreateBackToHueHomeAction()
	case hue_group_list.GroupSelect:
		m.selectedGroup = msg.Group
		m.details = hue_group_details.InitModel(m.adapter, msg.Group)
		m.state = details
	case hue_groups.BackToListAction:
		m.selectedGroup = nil
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

func (m *Model) init() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}
