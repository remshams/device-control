package home

import (
	"fmt"

	"github.com/remshams/device-control/tui/components/page_title"
	"github.com/remshams/device-control/tui/components/toast"
	"github.com/remshams/device-control/tui/lights/hue"
	"github.com/remshams/device-control/tui/lights/keylight"
	"github.com/remshams/device-control/tui/pages"
	hue_home_tabs "github.com/remshams/device-control/tui/pages/hue/home/tabs"
	keylight_home "github.com/remshams/device-control/tui/pages/keylight/home"
	"github.com/remshams/device-control/tui/stores"
	"github.com/remshams/device-control/tui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

const (
	menu      viewState = "menu"
	keylights viewState = "keylights"
	hueLights viewState = "hueLights"
)

type Model struct {
	keylight  keylight_home.Model
	hue       hue_home_tabs.Model
	menu      list.Model
	state     viewState
	toast     toast.Model
	pageTitle page_title.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter, hueAdapter *hue.HueAdapter) Model {
	return Model{
		keylight:  keylight_home.InitModel(keylightAdapter),
		hue:       hue_home_tabs.InitModel(hueAdapter),
		menu:      createMenu(),
		state:     menu,
		toast:     toast.InitModel(),
		pageTitle: page_title.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.toast, _ = m.toast.Update(msg)
	m.pageTitle, _ = m.pageTitle.Update(msg)
	if pages.IsSystemMsg(msg) {
		cmd = m.processSystemUpdate(msg)
	} else {
		switch m.state {
		case menu:
			cmd = m.processMenuUpdate(msg)
		case keylights:
			cmd = m.processKeylightsUpdate(msg)
		case hueLights:
			cmd = m.processHueUpate(msg)
		}
	}
	return m, cmd
}

func (m *Model) processSystemUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		stores.LayoutStore.Width = msg.Width
		stores.LayoutStore.Height = msg.Height
		cmd = pages.CreateWindowResizeAction(msg.Width, msg.Height)
	}
	return cmd
}

func (m *Model) processMenuUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.WindowResizeAction:
		h, v := styles.ListStyles.GetFrameSize()
		m.menu.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch m.state {
		case menu:
			switch msg.String() {
			case "ctrl+c", "q":
				cmd = tea.Quit
			case "enter":
				if m.menu.Index() == 0 {
					m.state = keylights
					cmd = m.keylight.Init()
				} else {
					m.state = hueLights
					cmd = m.hue.Init()
				}
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
	case pages.BackToMenuAction:
		m.state = menu
		cmd = page_title.CreateSetPageTitleMsg("")
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		default:
			m.keylight, cmd = m.keylight.Update(msg)
		}
	default:
		m.keylight, cmd = m.keylight.Update(msg)
	}
	return cmd
}

func (m *Model) processHueUpate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.BackToMenuAction:
		m.state = menu
		cmd = page_title.CreateSetPageTitleMsg("")
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		default:
			m.hue, cmd = m.hue.Update(msg)
		}
	default:
		m.hue, cmd = m.hue.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.renderPageTitle(), m.renderPageContent(), m.renderToast())
}

func (m Model) renderPageTitle() string {
	pageTitle := lipgloss.NewStyle().
		PaddingTop(styles.Padding).
		PaddingLeft(styles.Padding)
	return pageTitle.Render(m.pageTitle.View())
}

func (m Model) renderPageContent() string {
	switch m.state {
	case menu:
		return styles.ListStyles.Render(m.menu.View())
	case keylights:
		return m.keylight.View()
	case hueLights:
		return m.hue.View()
	default:
		return ""
	}
}

func (m Model) renderToast() string {
	toastStyle := lipgloss.NewStyle().PaddingTop(styles.Padding)
	return toastStyle.Render(m.toast.View())
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
