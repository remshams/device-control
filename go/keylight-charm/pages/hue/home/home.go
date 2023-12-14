package hue_home

import (
	"keylight-charm/lights/hue"
	"keylight-charm/pages"
	hue_groups "keylight-charm/pages/hue/groups"
	hue_groups_home "keylight-charm/pages/hue/groups/home"
	"keylight-charm/stores"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

var menuStyles = lipgloss.NewStyle().Margin(1, 2)

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
	menu   viewState = "menu"
	groups viewState = "groups"
)

type Model struct {
	adapter *hue.HueAdapter
	menu    list.Model
	groups  hue_groups_home.Model
	state   viewState
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		menu:    createMenu(),
		groups:  hue_groups_home.InitModel(adapter),
		state:   menu,
	}
}

func (m Model) Init() tea.Cmd {
	return m.createInitMsg()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.state == menu {
		cmd = m.processMenuUpdate(msg)
	} else {
		cmd = m.forwardUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processMenuUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case hue_groups.BackToGroupHomeAction:
		m.state = menu
	case initMsg:
		updateMenuLayout(&m.menu)
	case pages.WindowResizeAction:
		updateMenuLayout(&m.menu)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = pages.CreateBackToMenuAction()
		case "enter":
			if m.menu.Index() == 0 {
				m.state = groups
				cmd = m.groups.Init()
			}
		default:
			cmd = m.forwardUpdate(msg)
		}
	}
	return cmd
}

func updateMenuLayout(menu *list.Model) {
	h, v := menuStyles.GetFrameSize()
	menu.SetSize(stores.LayoutStore.Width-h, stores.LayoutStore.Height-v)
}

func (m *Model) forwardUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case menu:
		m.menu, cmd = m.menu.Update(msg)
	case groups:
		m.groups, cmd = m.groups.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case menu:
		return menuStyles.Render(m.menu.View())
	case groups:
		return m.groups.View()
	default:
		return ""
	}
}

func createMenu() list.Model {
	items := []list.Item{
		menuItem{title: "HueGroups", desc: "Control hue groups"},
		menuItem{title: "HueBridges", desc: "Manage hue bridges (pair...)"},
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
