package hue_home_tabs

import (
	"fmt"

	"github.com/remshams/device-control/tui/components/header"
	dc_tabs "github.com/remshams/device-control/tui/components/tabs"
	"github.com/remshams/device-control/tui/components/toast"
	"github.com/remshams/device-control/tui/lights/hue"
	"github.com/remshams/device-control/tui/pages"
	pages_hue "github.com/remshams/device-control/tui/pages/hue"
	hue_bridges_home "github.com/remshams/device-control/tui/pages/hue/bridges/home"
	hue_groups_home "github.com/remshams/device-control/tui/pages/hue/groups/home"
	hue_lights_home "github.com/remshams/device-control/tui/pages/hue/lights/home"
	"github.com/remshams/device-control/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState string

type initMsg struct{}

const (
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
		state:   groups,
	}
}

func (m Model) Init() tea.Cmd {
	return m.createInitMsg()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		cmd = tea.Batch(m.tabs.Init(), m.groups.Init(), header.CreateSetHeaderMsg("Hue Home"))
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
		cmd = pages.CreateBackToMenuAction()
	case dc_tabs.TabSelectedMsg:
		switch msg {
		case 0:
			m.state = groups
			cmd = m.groups.Init()
		case 1:
			m.state = bridges
			cmd = m.bridges.Init()
		case 2:
			m.state = lights
			cmd = m.lights.Init()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = m.forwardUpdate(msg)
		case "r":
			cmd = pages_hue.CreateReloadBridgesAction()
		default:
			cmd = m.forwardUpdate(msg)
		}

	default:
		cmd = m.forwardUpdate(msg)

	}
	return m, cmd
}

func (m *Model) forwardUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var tabsCmd tea.Cmd
	m.tabs, tabsCmd = m.tabs.Update(msg)
	switch m.state {
	case groups:
		m.groups, cmd = m.groups.Update(msg)
	case bridges:
		m.bridges, cmd = m.bridges.Update(msg)
	case lights:
		m.lights, cmd = m.lights.Update(msg)
	}
	return tea.Batch(tabsCmd, cmd)
}

func (m Model) View() string {
	body := ""
	switch m.state {
	case groups:
		body = m.groups.View()
	case bridges:
		body = m.bridges.View()
	case lights:
		body = m.lights.View()
	}
	return fmt.Sprintf(
		"\n%s\n%s",
		lipgloss.NewStyle().PaddingBottom(styles.Padding).Render(m.tabs.View()),
		body,
	)
}

func createTabs() dc_tabs.Model {
	return dc_tabs.New([]string{"Groups", "Bridges", "Lights"})
}

func (m Model) createInitMsg() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}
