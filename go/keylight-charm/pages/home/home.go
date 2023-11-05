package home

import (
	"keylight-charm/lights/keylight"
	keylight_home "keylight-charm/pages/keylight/home"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var menuStyle = lipgloss.NewStyle().Margin(1, 2)

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

const (
	menu      viewState = "menu"
	keylights viewState = "keylights"
	hue       viewState = "hue"
)

type Model struct {
	keylight keylight_home.Model
	menu     list.Model
	state    viewState
}

func InitModel(keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{
		keylight: keylight_home.InitModel(keylightAdapter),
		menu:     createMenu(),
		state:    menu,
	}
}

func (m Model) Init() tea.Cmd {
	return m.keylight.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case menu:
		cmd = m.processMenuUpdate(msg)
	case keylights:
		cmd = m.processKeylightsUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processMenuUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := menuStyle.GetFrameSize()
		m.menu.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch m.state {
		case menu:
			switch msg.String() {
			case "enter":
				m.state = keylights
				cmd = m.keylight.Init()
			default:
				m.menu, cmd = m.menu.Update(msg)
			}
		}
	default:
		m.menu, cmd = m.menu.Update(msg)
	}
	return cmd
}

func (m *Model) processKeylightsUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = menu
		default:
			m.keylight, cmd = m.keylight.Update(msg)
		}
	default:
		m.keylight, cmd = m.keylight.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case menu:
		return menuStyle.Render(m.menu.View())
	case keylights:
		return m.keylight.View()
	default:
		return ""
	}
}

func createMenu() list.Model {
	items := []list.Item{
		menuItem{title: "Keylights", desc: "Control keylights"},
		menuItem{title: "HueLights", desc: "Control huelights"},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Lights"
	return list
}
