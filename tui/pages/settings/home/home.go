package settings_home

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	device_control_settings "github.com/remshams/device-control/settings/public"
	"github.com/remshams/device-control/tui/components/page_title"
	dc_tabs "github.com/remshams/device-control/tui/components/tabs"
	settings_location "github.com/remshams/device-control/tui/pages/settings/location"
	"github.com/remshams/device-control/tui/styles"
)

type Model struct {
	settings *device_control_settings.Settings
	location settings_location.Model
	tabs     dc_tabs.Model
}

func InitModel(settings *device_control_settings.Settings) Model {
	return Model{
		settings: settings,
		location: settings_location.InitModel(settings),
		tabs:     dc_tabs.New([]string{"Location"}),
	}
}

func (m Model) Init() tea.Cmd {
	return page_title.CreateSetPageTitleMsg("Settings")
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.location, cmd = m.location.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		lipgloss.NewStyle().PaddingBottom(styles.Padding).Render(m.tabs.View()),
		m.location.View(),
	)
}
