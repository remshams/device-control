package settings_location

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	device_control_settings "github.com/remshams/device-control/settings/public"
	kl_cursor "github.com/remshams/device-control/tui/components/cursor"
	kl_textinput "github.com/remshams/device-control/tui/components/dc_textinput"
	"github.com/remshams/device-control/tui/components/page_help"
	page_settings "github.com/remshams/device-control/tui/pages/settings"
)

type keyMapNavMode struct {
	cursor    kl_cursor.KeyMap
	textinput kl_textinput.KeyMap
	Save      key.Binding
	Quit      key.Binding
}

func (k keyMapNavMode) ShortHelp() []key.Binding {
	return []key.Binding{
		k.cursor.Up,
		k.cursor.Down,
		k.textinput.Edit,
		k.Save,
		k.Quit,
	}
}

func (k keyMapNavMode) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var settingsLocationMap = keyMapNavMode{
	cursor:    kl_cursor.CursorKeyMap,
	textinput: kl_textinput.TextInputKeyMap,
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
	lat      kl_textinput.Model
	lng      kl_textinput.Model
	cursor   kl_cursor.CursorState
}

func InitModel(settings *device_control_settings.Settings) Model {
	m := Model{
		settings: settings,
		lat:      kl_textinput.New("Latitude", ""),
		lng:      kl_textinput.New("Longtitude", ""),
		cursor:   kl_cursor.InitCursorState(2),
	}
	m.lat.Input.SetValue(strconv.FormatFloat(settings.GetLatitude(), 'f', -1, 64))
	m.lng.Input.SetValue(strconv.FormatFloat(settings.GetLongtitude(), 'f', -1, 64))
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
			if m.cursor.Index() == 0 {
				m.lat, cmd = m.lat.Update(msg)
			} else {
				m.lng, cmd = m.lng.Update(msg)
			}
			m.cursor.Update(msg)
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		kl_cursor.RenderLine(
			m.lat.View(),
			m.cursor.Index() == 0,
			m.lat.Input.Focused(),
		),
		kl_cursor.RenderLine(
			m.lng.View(),
			m.cursor.Index() == 1,
			m.lng.Input.Focused()),
	)
}
