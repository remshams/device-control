package home

import (
	"keylight-charm/keylight"
	keylight_details "keylight-charm/pages/keylight/details"
	keylight_model "keylight-charm/pages/keylight/details/model"
	keylight_edit "keylight-charm/pages/keylight/edit"
	keylight_list "keylight-charm/pages/keylight/list"
	"keylight-control/control"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	initial viewState = "init"
	list              = "list"
	details           = "details"
	add               = "add"
)

type initMsg struct{}

type Model struct {
	keylightAdapter *keylight.KeylightAdapter
	state           viewState
	keylights       []control.Keylight
	list            keylight_list.Model
	details         *keylight_details.Model
	edit            *keylight_edit.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{keylightAdapter: keylightAdapter, state: initial, details: nil, edit: nil}
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
	case keylight_list.AddKeylight:
		newKeylight := keylight_edit.InitModel(nil)
		m.edit = &newKeylight
		m.state = add
	case keylight_model.AbortAction:
		m.details = nil
		m.edit = nil
		m.state = list
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		default:
			cmd = m.updateChilds(msg)
		}
	default:
		cmd = m.updateChilds(msg)
	}
	return m, cmd
}

func (m *Model) updateChilds(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case list:
		m.list, cmd = m.list.Update(msg)
	case details:
		var details keylight_details.Model
		details, cmd = m.details.Update(msg)
		m.details = &details
	case add:
		var edit keylight_edit.Model
		edit, cmd = m.edit.Update(msg)
		m.edit = &edit
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case initial:
		return "Loading..."
	case list:
		return m.list.View()
	case details:
		return m.details.View()
	case add:
		return m.edit.View()
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
