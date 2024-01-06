package hue_group_list

import (
	"fmt"
	"strconv"

	hue_control "github.com/remshams/device-control/hue-control/pubilc"
	"github.com/remshams/device-control/tui/components/page_help"
	kl_table "github.com/remshams/device-control/tui/components/table"
	"github.com/remshams/device-control/tui/components/toast"
	"github.com/remshams/device-control/tui/lights/hue"
	pages_hue "github.com/remshams/device-control/tui/pages/hue"
	hue_groups "github.com/remshams/device-control/tui/pages/hue/groups"
	hue_group_details "github.com/remshams/device-control/tui/pages/hue/groups/details"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Quit                 key.Binding
	ToggleAllGroupLights key.Binding
	SelectGroup          key.Binding
	ReloadBridges        key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.ToggleAllGroupLights, k.SelectGroup, k.ReloadBridges}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var defaultKeyMap = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("esc", "esc"),
		key.WithHelp("esc", "Go back"),
	),
	ToggleAllGroupLights: key.NewBinding(
		key.WithKeys("t", "t"),
		key.WithHelp("t", "Toggle all group lights"),
	),
	SelectGroup: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select group"),
	),
	ReloadBridges: key.NewBinding(
		key.WithKeys("r", "r"),
		key.WithHelp("r", "Reload groups"),
	),
}

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

func (m Model) Init() tea.Cmd {
	return page_help.CreateSetKeyMapMsg(defaultKeyMap)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages_hue.BridgesReloadedAction:
		m.table.SetRows(createTableRows(m.adapter.Control.GetBridges()))
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultKeyMap.ToggleAllGroupLights):
			cmd = m.toggleAllGroupLights()
		case key.Matches(msg, defaultKeyMap.SelectGroup):
			cmd = m.selectGroup(m.table.SelectedRow()[0])
		case key.Matches(msg, defaultKeyMap.Quit):
			cmd = hue_groups.CreateBackToGroupHomeAction()
		case key.Matches(msg, defaultKeyMap.ReloadBridges):
			cmd = pages_hue.CreateReloadBridgesAction()
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
