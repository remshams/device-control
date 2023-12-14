package hue_bridges_list

import (
	hue_control "hue-control/pubilc"
	kl_table "keylight-charm/components/table"
	"keylight-charm/lights/hue"
	hue_bridges "keylight-charm/pages/hue/bridges"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	adapter           *hue.HueAdapter
	discoveredBridges []hue_control.DiscoveredBridge
	bridges           []hue_control.Bridge
	table             table.Model
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter:           adapter,
		discoveredBridges: adapter.Control.GetDiscoveredBridges(),
		bridges:           adapter.Control.GetBridges(),
		table: createTable(
			adapter.Control.GetBridges(),
			adapter.Control.GetDiscoveredBridges(),
		),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = hue_bridges.CreateBackToBridgesHomeAction()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return "Bridges"
}

func createTable(bridges []hue_control.Bridge, discoveredBridges []hue_control.DiscoveredBridge) table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 5},
		{Title: "Ip", Width: 15},
		{Title: "ApiKey", Width: 15},
	}
	rows := []table.Row{}

	for _, bridge := range bridges {
		rows = append(rows, table.Row{
			bridge.GetId(),
			bridge.GetIp().String(),
			bridge.GetApiKey(),
		})
	}

	for _, bridge := range discoveredBridges {
		rows = append(rows, table.Row{
			bridge.Id,
			bridge.Ip.String(),
			"",
		})
	}

	return kl_table.CreateTable(columns, rows)
}
