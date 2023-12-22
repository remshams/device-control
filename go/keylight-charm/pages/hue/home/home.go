package hue_home

import (
	"keylight-charm/lights/hue"
	"keylight-charm/pages"
	pages_hue "keylight-charm/pages/hue"
	hue_bridges_home "keylight-charm/pages/hue/bridges/home"
	hue_groups_home "keylight-charm/pages/hue/groups/home"
	hue_lights_home "keylight-charm/pages/hue/lights/home"
	"keylight-charm/stores"
	"keylight-charm/styles"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

type menuItem struct {
	title, desc string
}

func (menuItem menuItem) Title() string {
	return menuItem.title
}

func (menuItem menuItem) Description() string {
	return menuItem.desc
}

func (menuItem menuItem) FilterValue() string { return menuItem.title }

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
	menu    list.Model
	bridges hue_bridges_home.Model
	groups  hue_groups_home.Model
	lights  hue_lights_home.Model
	state   viewState
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		menu:    createMenu(),
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
		m.adapter.Control.LoadBridges()
		cmd = pages_hue.CreateBridgesReloadedAction()
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
	case initMsg:
		updateMenuLayout(&m.menu)
	case pages.WindowResizeAction:
		updateMenuLayout(&m.menu)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = pages.CreateBackToMenuAction()
		case "enter":
			switch m.menu.Index() {
			case 0:
				m.state = bridges
				cmd = m.bridges.Init()
			case 1:
				m.state = groups
				cmd = m.groups.Init()
			case 2:
				m.state = lights
				cmd = m.lights.Init()
			}
		default:
			cmd = m.forwardUpdate(msg)
		}
	}
	return cmd
}

func updateMenuLayout(menu *list.Model) {
	h, v := styles.ListStyles.GetFrameSize()
	menu.SetSize(stores.LayoutStore.Width-h, stores.LayoutStore.Height-v)
}

func (m *Model) forwardUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case menu:
		m.menu, cmd = m.menu.Update(msg)
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
		return styles.ListStyles.Render(m.menu.View())
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

func createMenu() list.Model {
	items := []list.Item{
		menuItem{title: "Bridges", desc: "Manage hue bridges (pair...)"},
		menuItem{title: "Groups", desc: "Control hue groups"},
		menuItem{title: "Lights", desc: "Control hue lights"},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Hue Home"
	return list
}

func (m Model) createInitMsg() tea.Cmd {
	return func() tea.Msg {
		return initMsg{}
	}
}
