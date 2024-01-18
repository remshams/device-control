package settings_location

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	device_control_settings "github.com/remshams/device-control/settings/public"
	kl_cursor "github.com/remshams/device-control/tui/components/cursor"
	"github.com/remshams/device-control/tui/components/page_help"
	kl_textinput "github.com/remshams/device-control/tui/components/textinput"
	page_settings "github.com/remshams/device-control/tui/pages/settings"
)

type keyMap struct {
	cursor kl_cursor.KeyMap
	Save   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.cursor.Up,
		k.cursor.Down,
		k.Save,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var settingsLocationMap = keyMap{
	cursor: kl_cursor.CursorKeyMap,
	Save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Go back"),
	),
}

type Model struct {
	settings *device_control_settings.Settings
	lat      textinput.Model
	lng      textinput.Model
	cursor   kl_cursor.CursorState
}

func InitModel(settings *device_control_settings.Settings) Model {
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
	return page_help.CreateSetKeyMapMsg(settingsLocationMap)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, settingsLocationMap.Quit):
			cmd = page_settings.CreateBackToSettingsHomeAction()
		default:
			m.cursor.Update(msg)
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		kl_cursor.RenderLine(
			kl_textinput.CreateTextInputView(m.lat, "Latitude", ""),
			m.cursor.Index() == 0,
			m.lat.Focused(),
		),
		kl_cursor.RenderLine(
			kl_textinput.CreateTextInputView(m.lng, "Longtitude", ""),
			m.cursor.Index() == 1,
			m.lng.Focused()),
	)
}
