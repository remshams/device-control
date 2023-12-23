package hue_group_list

import (
	"fmt"
	hue_control "hue-control/pubilc"
	kl_table "keylight-charm/components/table"
	"keylight-charm/components/toast"
	"keylight-charm/lights/hue"
	pages_hue "keylight-charm/pages/hue"
	hue_groups "keylight-charm/pages/hue/groups"
	hue_group_details "keylight-charm/pages/hue/groups/details"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type GroupSelect struct {
	Group *hue_control.Group
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
		table:   kl_table.CreateTable(createTableColumns(), createTableRows(adapter.Control.GetBridges())),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages_hue.BridgesReloadedAction:
		m.table.SetRows(createTableRows(m.adapter.Control.GetBridges()))
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

func createTableColumns() []table.Column {
	return []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Number of lights", Width: 20},
		{Title: "Bridge Id", Width: 20},
		{Title: "Lights On", Width: 10},
	}
}

func createTableRows(bridges []hue_control.Bridge) []table.Row {
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
					strconv.FormatBool(group.GetOn()),
				},
			)
		}
	}

	return rows
}

func (m *Model) selectGroup(id string) tea.Cmd {
	return func() tea.Msg {
		return GroupSelect{
			Group: m.findSelectedGroup(id),
		}
	}
}

func (m Model) toggleAllGroupLights() tea.Cmd {
	group := m.findSelectedGroup(m.table.SelectedRow()[0])
	on := !group.GetOn()
	group.SetOn(on)
	err := group.SendGroup()
	if err != nil {
		return toast.CreateErrorToastAction("Failed to toggle group")
	}
	onText := "on"
	if on {
		onText = "off"
	}
	return tea.Batch(
		toast.CreateSuccessToastAction(fmt.Sprintf("All groups set to %s", onText)),
		pages_hue.CreateBridgesReloadedAction(),
	)
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
