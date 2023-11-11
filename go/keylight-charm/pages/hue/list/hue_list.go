package hue_list

import (
	hue_control "hue-control/pubilc"
	kl_table "keylight-charm/components/table"
	"keylight-charm/lights/hue"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	adapter *hue.HueAdapter
	bridges []hue_control.Bridge
	table   table.Model
}

func InitModel(adapter *hue.HueAdapter, bridges []hue_control.Bridge) Model {
	return Model{
		adapter: adapter,
		bridges: bridges,
		table:   createTable(bridges),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return m.table.View()
}

func createTable(bridges []hue_control.Bridge) table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Number of lights", Width: 20},
		{Title: "Connected", Width: 10},
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
					strconv.FormatBool(group.GetConnected()),
				},
			)
		}
	}
	return kl_table.CreateTable(columns, rows)
}
