package hue_bridges_list

import (
	kl_table "keylight-charm/components/table"
	"keylight-charm/components/toast"
	"keylight-charm/lights/hue"
	pages_hue "keylight-charm/pages/hue"
	hue_bridges "keylight-charm/pages/hue/bridges"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type initMsg struct{}
type bridgesDiscovered struct{}
type bridgePaired struct {
	success bool
	message string
}

type viewState string

const (
	initial viewState = "initial"
	list    viewState = "list"
)

type Model struct {
	adapter *hue.HueAdapter
	table   table.Model
	state   viewState
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
	case pages_hue.BridgesReloadedAction:
		m.table.SetRows(m.createTableRows())
	case bridgesDiscovered:
		cmd = tea.Batch(toast.CreateInfoToastAction("Bridges discovered"), m.reloadBridges())
	case bridgePaired:
		var toastCmd tea.Cmd
		if msg.success {
			toastCmd = toast.CreateSuccessToastAction(msg.message)
		} else {
			toastCmd = toast.CreateErrorToastAction(msg.message)
		}
		cmd = tea.Batch(toastCmd, m.reloadBridges())
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			cmd = tea.Batch(toast.CreateInfoToastAction("Discovering bridges..."), m.createBridgesDiscoveredMsg())
		case "p":
			cmd = tea.Batch(toast.CreateInfoToastAction("Pairing bridge, please press the button..."), m.createBridgePairedMsg())
		case "esc":
			cmd = hue_bridges.CreateBackToBridgesHomeAction()
		default:
			m.table, cmd = m.table.Update(msg)
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

	for _, bridge := range m.adapter.Control.GetNewlyDiscoveredBridges() {
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
	return func() tea.Msg {
		return pages_hue.BridgesReloadedAction{
			Bridges: m.adapter.Control.GetBridges(),
		}
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

func (m Model) createBridgePairedMsg() tea.Cmd {
	return func() tea.Msg {
		bridgeId := m.table.SelectedRow()[0]
		_, err := m.adapter.Control.Pair(bridgeId)
		if err != nil {
			return bridgePaired{success: false, message: err.Error()}
		} else {
			return bridgePaired{success: true, message: "Bridge paired"}
		}
	}
}
