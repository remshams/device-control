package hue_bridges_list

import (
	"fmt"

	"github.com/remshams/device-control/tui/components/page_help"
	kl_table "github.com/remshams/device-control/tui/components/table"
	"github.com/remshams/device-control/tui/components/toast"
	"github.com/remshams/device-control/tui/lights/hue"
	pages_hue "github.com/remshams/device-control/tui/pages/hue"
	hue_bridges "github.com/remshams/device-control/tui/pages/hue/bridges"

	"github.com/charmbracelet/bubbles/key"
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
	initial      viewState = "initial"
	list         viewState = "list"
	deleteBridge viewState = "deleteBridge"
)

type keyMap struct {
	Quit            key.Binding
	DiscoverBridges key.Binding
	PairBridge      key.Binding
	DeleteBridge    key.Binding
	ReloadBridges   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.DiscoverBridges, k.PairBridge, k.DeleteBridge, k.ReloadBridges}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var defaultKeyMap = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("esc", "esc"),
		key.WithHelp("esc", "Go back"),
	),
	DiscoverBridges: key.NewBinding(
		key.WithKeys("f", "f"),
		key.WithHelp("f", "Discover bridges"),
	),
	PairBridge: key.NewBinding(
		key.WithKeys("p", "p"),
		key.WithHelp("p", "Pair bridge"),
	),
	DeleteBridge: key.NewBinding(
		key.WithKeys("d", "d"),
		key.WithHelp("d", "Delete bridge"),
	),
	ReloadBridges: key.NewBinding(
		key.WithKeys("r", "r"),
		key.WithHelp("r", "Reload bridges"),
	),
}

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
	return tea.Batch(page_help.CreateSetKeyMapMsg(defaultKeyMap), m.createInitMsg())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		cmd = m.reloadBridges()
		m.table.SetRows(m.createTableRows())
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
		switch {
		case key.Matches(msg, defaultKeyMap.DiscoverBridges):
			cmd = tea.Batch(toast.CreateInfoToastAction("Discovering bridges..."), m.createBridgesDiscoveredMsg())
		case key.Matches(msg, defaultKeyMap.PairBridge):
			cmd = tea.Batch(toast.CreateInfoToastAction("Pairing bridge, please press the button..."), m.createBridgePairedMsg())
		case key.Matches(msg, defaultKeyMap.Quit):
			cmd = hue_bridges.CreateBackToBridgesHomeAction()
		case key.Matches(msg, defaultKeyMap.DeleteBridge):
			m.state = deleteBridge
			cmd = toast.CreateWarningToastAction(fmt.Sprintf("Are you sure you want to delete bridge %s? (y/n)", m.table.SelectedRow()[0]))
		case msg.String() == "y":
			if m.state == deleteBridge {
				cmd = m.deleteBridge()
				m.state = list
			}
		case key.Matches(msg, defaultKeyMap.ReloadBridges):
			cmd = pages_hue.CreateReloadBridgesAction()
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m *Model) deleteBridge() tea.Cmd {
	bridgeId := m.table.SelectedRow()[0]
	err := m.adapter.Control.RemoveBridge(bridgeId)
	if err != nil {
		return toast.CreateErrorToastAction("Could not delete bridge")
	} else {
		m.table.MoveUp(1)
		return tea.Batch(toast.CreateSuccessToastAction(fmt.Sprintf("Bridge %s deleted", bridgeId)), m.reloadBridges())
	}
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
	return pages_hue.CreateBridgesReloadedAction()
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
