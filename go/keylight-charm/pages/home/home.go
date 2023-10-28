package home

import (
	"fmt"
	"keylight-charm/components/actions"
	"keylight-charm/components/toast"
	"keylight-charm/keylight"
	"keylight-charm/pages/keylight/details"
	"keylight-charm/pages/keylight/edit"
	"keylight-charm/pages/keylight/list"
	"keylight-control/control"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	initial viewState = "init"
	list              = "list"
	details           = "details"
	add               = "add"
	edit              = "edit"
)

type initMsg struct{}

type reloadKeylights struct{}

type Model struct {
	keylightAdapter *keylight.KeylightAdapter
	state           viewState
	keylights       []control.Keylight
	list            keylight_list.Model
	details         *keylight_details.Model
	edit            *keylight_edit.Model
	toast           toast.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{
		keylightAdapter: keylightAdapter,
		state:           initial,
		details:         nil,
		edit:            nil,
		toast:           toast.InitModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.toast, _ = m.toast.Update(msg)
	switch msg := msg.(type) {
	case initMsg:
		m.keylights = m.keylightAdapter.Control.Keylights()
		m.list = keylight_list.InitModel(m.keylightAdapter, m.keylights)
		m.state = list
	case reloadKeylights:
		m.keylights = m.keylightAdapter.Control.Keylights()
		m.list = keylight_list.InitModel(m.keylightAdapter, m.keylights)
	case keylight_list.SelectedKeylight:
		keylightDetails := keylight_details.InitModel(msg.Keylight, m.keylightAdapter)
		m.details = &keylightDetails
		m.state = details
	case keylight_list.AddKeylight:
		newKeylight := keylight_edit.InitModel(nil, m.keylightAdapter)
		m.edit = &newKeylight
		m.state = add
	case keylight_list.EditKeylight:
		editKeylight := keylight_edit.InitModel(msg.Keylight, m.keylightAdapter)
		m.edit = &editKeylight
		m.state = edit
	case keylight_list.RemoveKeylight:
		_, err := m.keylightAdapter.RemoveKeylight(msg.Keylight.Metadata.Id)
		if err != nil {
			cmd = toast.CreateErrorToastAction("Keylight could not be deleted")
		} else {
			cmd = tea.Batch(toast.CreateInfoToastAction("Keylight deleted"), m.reloadKeylights())
		}
	case keylight_list.ReloadKeylights:
		cmd = m.reloadKeylights()
	case actions.SaveAction:
		cmd = m.reloadKeylights()
	case actions.AbortAction:
		m.details = nil
		m.edit = nil
		m.state = list
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		case "r":
			cmd = tea.Batch(m.reloadKeylights(), toast.CreateInfoToastAction("Keylight list reloaded"))
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
	case edit:
		var edit keylight_edit.Model
		edit, cmd = m.edit.Update(msg)
		m.edit = &edit
	}
	return cmd
}

func (m Model) View() string {
	component := ""
	switch m.state {
	case initial:
		return "Loading..."
	case list:
		component = m.list.View()
	case details:
		component = m.details.View()
	case add:
		component = m.edit.View()
	case edit:
		component = m.edit.View()
	default:
		return "Error"
	}

	return fmt.Sprintf("%s\n%s", component, m.toast.View())
}

func (m *Model) init() tea.Cmd {
	return func() tea.Msg {
		m.keylightAdapter.Control.LoadOrDiscoverKeylights()
		return initMsg{}
	}
}

func (m *Model) reloadKeylights() tea.Cmd {
	return func() tea.Msg {
		m.keylightAdapter.Control.LoadOrDiscoverKeylights()
		return reloadKeylights{}
	}
}
