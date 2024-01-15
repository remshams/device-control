package settings_home

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	device_control_settings "github.com/remshams/device-control/settings/public"
	kl_cursor "github.com/remshams/device-control/tui/components/cursor"
)

type keyMap struct {
	Save key.Binding
	Quit key.Binding
}

var settingsHomeKeyMap = keyMap{
	Save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save"),
	),
}

type Model struct {
	settings *device_control_settings.Settings
	lat      textinput.Model
	lng      textinput.Model
	cursor   kl_cursor.CursorState
}

func InitMolde(settings *device_control_settings.Settings) Model {
	m := Model{
		settings: settings,
		lat:      textinput.New(),
		lng:      textinput.New(),
		cursor:   kl_cursor.InitCursorState(2),
	}
	m.lat.SetValue(strconv.FormatFloat(settings.GetLatitude(), 'f', -1, 64))
	m.lng.SetValue(strconv.FormatFloat(settings.GetLongtitude(), 'f', -1, 64))
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		kl_cursor.RenderLine(
			m.lat.View(),
			m.cursor.Index() == 0,
			m.lat.Focused(),
		),
		kl_cursor.RenderLine(
			m.lng.View(),
			m.cursor.Index() == 1,
			m.lng.Focused()),
	)
}
