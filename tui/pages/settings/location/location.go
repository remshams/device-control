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

type viewState string

const (
	navigate viewState = "navigate"
	edit     viewState = "edit"
)

type Model struct {
	settings *device_control_settings.Settings
	lat      kl_textinput.Model
	lng      kl_textinput.Model
	cursor   kl_cursor.CursorState
	state    viewState
}

func InitModel(settings *device_control_settings.Settings) Model {
	m := Model{
		settings: settings,
		lat:      kl_textinput.New("Latitude", ""),
		lng:      kl_textinput.New("Longtitude", ""),
		cursor:   kl_cursor.InitCursorState(2),
		state:    navigate,
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
			if m.state == navigate {
				cmd = page_settings.CreateBackToSettingsHomeAction()
			} else {
				m.state = navigate
				cmd = m.resetSettings(msg)
			}
		case key.Matches(msg, settingsLocationMap.textinput.Apply):
			if m.state == edit {
				cmd = m.updateSelectedInput(msg)
				m.state = navigate
			}
		case key.Matches(msg, settingsLocationMap.textinput.Edit):
			if m.state == navigate {
				m.state = edit
				cmd = m.updateSelectedInput(msg)
			} else {
				cmd = m.updateSelectedInput(msg)
			}
		case key.Matches(msg, settingsLocationMap.Save):
			if m.state == navigate {
				m.state = navigate
			} else {
				cmd = m.updateSelectedInput(msg)
			}
			// TODO: save settings
		default:
			if m.state == navigate {
				m.cursor.Update(msg)
			} else {
				cmd = m.updateSelectedInput(msg)
			}
		}
	}
	return m, cmd
}

func (m *Model) updateSelectedInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.cursor.Index() == 0 {
		m.lat, cmd = m.lat.Update(msg)
	} else {
		m.lng, cmd = m.lng.Update(msg)
	}
	return cmd
}

func (m *Model) resetSettings(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.cursor.Index() == 0 {
		m.lat, cmd = m.lat.Update(msg)
		m.lat.Input.SetValue(strconv.FormatFloat(m.settings.GetLatitude(), 'f', -1, 64))
	} else {
		m.lng, cmd = m.lng.Update(msg)
		m.lng.Input.SetValue(strconv.FormatFloat(m.settings.GetLongtitude(), 'f', -1, 64))
	}
	return cmd
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
