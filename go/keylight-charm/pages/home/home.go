package home

import (
	"keylight-charm/keylight"
	keylight_details "keylight-charm/pages/keylight/details"
	keylight_list "keylight-charm/pages/keylight/list"
	"keylight-control/control"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	initial viewState = "init"
	list              = "list"
	details           = "details"
)

type initMsg struct{}

type Model struct {
	keylightAdapter *keylight.KeylightAdapter
	state           viewState
	keylights       []control.Keylight
	list            keylight_list.Model
	details         *keylight_details.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{keylightAdapter: keylightAdapter, state: initial, details: nil}
}

func (m Model) Init() tea.Cmd {
	return m.discoverKeylights()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		m.keylights = m.keylightAdapter.Control.Keylights()
		m.list = keylight_list.InitModel(m.keylightAdapter, m.keylights)
		m.state = list
	case keylight_list.SelectedKeylight:
		keylightDetails := keylight_details.InitModel(msg.Keylight, m.keylightAdapter)
		m.details = &keylightDetails
		m.state = details
	case keylight_details.AbortAction:
		m.details = nil
		m.state = list
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		default:
			switch m.state {
			case list:
				m.list, cmd = m.list.Update(msg)
			case details:
				var details keylight_details.Model
				details, cmd = m.details.Update(msg)
				m.details = &details
			}
		}
	default:
		switch m.state {
		case list:
			m.list, cmd = m.list.Update(msg)
		case details:
			var details keylight_details.Model
			details, cmd = m.details.Update(msg)
			m.details = &details
		}
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case initial:
		return "Loading..."
	case list:
		return m.list.View()
	case details:
		return m.details.View()
	default:
		return "Error"
	}
}

func (m *Model) discoverKeylights() tea.Cmd {
	return func() tea.Msg {
		m.keylightAdapter.Control.LoadOrDiscoverKeylights()
		return initMsg{}
	}
}
