package hue_home_tabs

import (
	dc_tabs "github.com/remshams/device-control/tui/components/tabs"
	"github.com/remshams/device-control/tui/components/toast"
	"github.com/remshams/device-control/tui/lights/hue"
	"github.com/remshams/device-control/tui/pages"
	pages_hue "github.com/remshams/device-control/tui/pages/hue"
	hue_bridges_home "github.com/remshams/device-control/tui/pages/hue/bridges/home"
	hue_groups_home "github.com/remshams/device-control/tui/pages/hue/groups/home"
	hue_lights_home "github.com/remshams/device-control/tui/pages/hue/lights/home"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

type initMsg struct{}

const (
	menu    viewState = "menu"
	groups  viewState = "groups"
	bridges viewState = "bridges"
	lights  viewState = "lights"
)

type Model struct {
	adapter *hue.HueAdapter
	tabs    dc_tabs.Model
	bridges hue_bridges_home.Model
	groups  hue_groups_home.Model
	lights  hue_lights_home.Model
	state   viewState
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		tabs:    createTabs(),
		bridges: hue_bridges_home.InitModel(adapter),
		groups:  hue_groups_home.InitModel(adapter),
		lights:  hue_lights_home.InitModel(adapter),
		state:   menu,
	}
}

func (m Model) Init() tea.Cmd {
	return m.createInitMsg()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages_hue.ReloadBridgesAction:
		err := m.adapter.Control.LoadBridges()
		if err != nil {
			cmd = toast.CreateErrorToastAction("Failed to reload bridges")
		}
		cmd = tea.Batch(
			toast.CreateSuccessToastAction("Bridges/Groups/Lights reloaded"),
			pages_hue.CreateBridgesReloadedAction(),
		)
	case pages_hue.BackToHueHomeAction:
		m.state = menu
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			cmd = pages_hue.CreateReloadBridgesAction()
		default:
			cmd = m.defaultUpdate(msg)
		}

	default:
		cmd = m.defaultUpdate(msg)

	}
	return m, cmd
}

func (m *Model) defaultUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.state == menu {
		cmd = m.processMenuUpdate(msg)
	} else {
		cmd = m.forwardUpdate(msg)
	}
	return cmd
}

func (m *Model) processMenuUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = pages.CreateBackToMenuAction()
		default:
			cmd = m.forwardUpdate(msg)
		}
	}
	return cmd
}

func (m *Model) forwardUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case menu:
		m.tabs, cmd = m.tabs.Update(msg)
	case groups:
		m.groups, cmd = m.groups.Update(msg)
	case bridges:
		m.bridges, cmd = m.bridges.Update(msg)
	case lights:
		m.lights, cmd = m.lights.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case menu:
		return m.tabs.View()
	case groups:
		return m.groups.View()
	case bridges:
		return m.bridges.View()
	case lights:
		return m.lights.View()
	default:
		return ""
	}
}

func createTabs() dc_tabs.Model {
	return dc_tabs.New([]string{"Bridges", "Groups", "Lights"})
}

func (m Model) createInitMsg() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}
