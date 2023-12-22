package hue_lights_list

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
	lights  table.Model
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		lights: kl_table.CreateTable(
			createLightsColumns(),
			createLightsRows(adapter.Control.GetBridges()),
		),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.lights, cmd = m.lights.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.lights.View()
}

func createLightsColumns() []table.Column {
	return []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "On", Width: 5},
		{Title: "BridgeId", Width: 20},
	}
}

func createLightsRows(bridges []hue_control.Bridge) []table.Row {
	var rows []table.Row

	for _, bridge := range bridges {
		lights := bridge.GetLights()
		for _, light := range lights {
			rows = append(rows, table.Row{
				light.GetId(),
				light.GetName(),
				strconv.FormatBool(light.GetOn()),
				light.GetBridgeId(),
			})
		}
	}
	return rows
}
