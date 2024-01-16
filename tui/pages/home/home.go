package home

import (
	"fmt"

	device_control_settings "github.com/remshams/device-control/settings/public"
	"github.com/remshams/device-control/tui/components/page_help"
	"github.com/remshams/device-control/tui/components/page_title"
	"github.com/remshams/device-control/tui/components/toast"
	"github.com/remshams/device-control/tui/lights/hue"
	"github.com/remshams/device-control/tui/lights/keylight"
	"github.com/remshams/device-control/tui/pages"
	"github.com/remshams/device-control/tui/pages/hue/home"
	keylight_home "github.com/remshams/device-control/tui/pages/keylight/home"
	settings_home "github.com/remshams/device-control/tui/pages/settings/home"
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
	settings  viewState = "settings"
)

type Model struct {
	keylight     keylight_home.Model
	hue          hue_home_tabs.Model
	settingsPage settings_home.Model
	menu         list.Model
	state        viewState
	toast        toast.Model
	keyMap       page_help.Model
	pageTitle    page_title.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter, hueAdapter *hue.HueAdapter, settings *device_control_settings.Settings) Model {
	return Model{
		keylight:     keylight_home.InitModel(keylightAdapter),
		hue:          hue_home_tabs.InitModel(hueAdapter),
		settingsPage: settings_home.InitModel(settings),
		menu:         createMenu(),
		state:        menu,
		toast:        toast.InitModel(),
		keyMap:       page_help.New(),
		pageTitle:    page_title.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.toast, _ = m.toast.Update(msg)
	m.pageTitle, _ = m.pageTitle.Update(msg)
	m.keyMap, _ = m.keyMap.Update(msg)
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
		case settings:
			cmd = m.processSettingsUpdate(msg)
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
				switch m.menu.Index() {
				case 0:
					m.state = keylights
					cmd = m.keylight.Init()
				case 1:
					m.state = hueLights
					cmd = m.hue.Init()
				case 2:
					m.state = settings
					cmd = m.settingsPage.Init()
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
		cmd = m.resetLayout()
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
		cmd = m.resetLayout()
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

func (m *Model) processSettingsUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.BackToMenuAction:
		m.state = menu
		cmd = m.resetLayout()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		default:
			m.settingsPage, cmd = m.settingsPage.Update(msg)
		}
	default:
		m.settingsPage, cmd = m.settingsPage.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.renderPageTitle(),
		m.renderPageContent(),
		m.renderToast(),
		m.keyMap.View(),
	)
}

func (m Model) renderPageTitle() string {
	pageTitle := lipgloss.NewStyle().
		PaddingTop(styles.Padding).
		PaddingBottom(styles.Padding).
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
	case settings:
		return m.settingsPage.View()
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
		menuItem{title: "Settings", desc: "Settings"},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Devices"
	return list
}

func (m Model) resetLayout() tea.Cmd {
	return tea.Batch(page_title.CreateSetPageTitleMsg(""), page_help.CreateResetKeyMapMsg())
}
