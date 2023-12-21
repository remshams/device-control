package hue_group_list

import (
	hue_control "hue-control/pubilc"
	kl_table "keylight-charm/components/table"
	"keylight-charm/lights/hue"
	pages_hue "keylight-charm/pages/hue"
	hue_groups "keylight-charm/pages/hue/groups"
	hue_group_details "keylight-charm/pages/hue/groups/details"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type GroupSelect struct {
	Group hue_control.Group
}

type Model struct {
	adapter *hue.HueAdapter
	table   table.Model
	details hue_group_details.Model
}

func InitModel(adapter *hue.HueAdapter) Model {
	adapter.Control.LoadBridges()
	return Model{
		adapter: adapter,
		table:   createTable(adapter.Control.GetBridges()),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			cmd = m.toggleAllGroupLights()
		case "enter":
			cmd = m.selectGroup(m.table.SelectedRow()[0])
		case "esc":
			cmd = hue_groups.CreateBackToGroupHomeAction()
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

func createTable(bridges []hue_control.Bridge) table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Number of lights", Width: 20},
		{Title: "Bridge Id", Width: 20},
	}
	rows := []table.Row{}

	for _, bridge := range bridges {
		for _, group := range bridge.GetGroups() {
			rows = append(
				rows,
				table.Row{
					group.GetId(),
					group.GetName(),
					strconv.Itoa(len(group.GetLightIds())),
					group.GetBridgeId(),
				},
			)
		}
	}
	return kl_table.CreateTable(columns, rows)
}

func (m *Model) selectGroup(id string) tea.Cmd {
	return func() tea.Msg {
		return GroupSelect{
			Group: *m.findSelectedGroup(id),
		}
	}
}

func (m Model) toggleAllGroupLights() tea.Cmd {
	group := m.findSelectedGroup(m.table.SelectedRow()[0])
	group.SetOn(!group.GetOn())
	group.SendGroup()
	m.adapter.Control.LoadBridges()
	return pages_hue.CreateBridgesReloadedAction()
}

func (m Model) findSelectedGroup(id string) *hue_control.Group {
	for _, bridge := range m.adapter.Control.GetBridges() {
		group := bridge.GetGroupById(id)
		if group != nil {
			return group
		}
	}
	return nil
}
