package hue_bridges_list

import (
	hue_control "hue-control/pubilc"
	kl_table "keylight-charm/components/table"
	"keylight-charm/components/toast"
	"keylight-charm/lights/hue"
	hue_bridges "keylight-charm/pages/hue/bridges"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type initMsg struct{}
type bridgesDiscovered struct{}
type bridgesReloaded struct{}

type viewState string

const (
	initial viewState = "initial"
	list    viewState = "list"
)

type Model struct {
	adapter           *hue.HueAdapter
	discoveredBridges []hue_control.DiscoveredBridge
	bridges           []hue_control.Bridge
	table             table.Model
	state             viewState
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		state:   initial,
		table:   createInitialTable(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.createInitMsg()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		cmd = m.reloadBridges()
		m.state = list
	case bridgesReloaded:
		m.table.SetRows(m.createTableRows())
	case bridgesDiscovered:
		cmd = tea.Batch(toast.CreateInfoToastAction("Bridges discovered"), m.reloadBridges())
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			cmd = tea.Batch(toast.CreateInfoToastAction("Discovering bridges..."), m.createBridgesDiscoveredMsg())
		case "esc":
			cmd = hue_bridges.CreateBackToBridgesHomeAction()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case initial:
		return "Loading..."
	case list:
		return m.table.View()
	default:
		return ""
	}
}

func (m Model) createTable(columns []table.Column, rows []table.Row) table.Model {
	return kl_table.CreateTable(columns, rows)
}

func createInitialTable() table.Model {
	return kl_table.CreateTable(createTableColumns(), []table.Row{})
}

func createTableColumns() []table.Column {
	return []table.Column{
		{Title: "Id", Width: 40},
		{Title: "Ip", Width: 15},
		{Title: "ApiKey", Width: 45},
	}
}

func (m Model) createTableRows() []table.Row {
	rows := []table.Row{}

	for _, bridge := range m.adapter.Control.GetBridges() {
		rows = append(rows, table.Row{
			bridge.GetId(),
			bridge.GetIp().String(),
			bridge.GetApiKey(),
		})
	}

	for _, bridge := range m.adapter.Control.GetDiscoveredBridges() {
		rows = append(rows, table.Row{
			bridge.Id,
			bridge.Ip.String(),
			"",
		})
	}

	return rows
}

func (m *Model) reloadBridges() tea.Cmd {
	m.adapter.Control.LoadBridges()
	m.bridges = m.adapter.Control.GetBridges()
	m.discoveredBridges = m.adapter.Control.GetDiscoveredBridges()
	return func() tea.Msg {
		return bridgesReloaded{}
	}
}

func (m Model) createBridgesDiscoveredMsg() tea.Cmd {
	return func() tea.Msg {
		m.adapter.Control.DiscoverBridges()
		return bridgesDiscovered{}
	}
}

func (m Model) createInitMsg() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}
