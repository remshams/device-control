package keylight_home

import (
	hue_control "github.com/remshams/device-control/hue-control/pubilc"
	"github.com/remshams/device-control/tui/components/toast"
	"github.com/remshams/device-control/tui/lights/keylight"
	pages_keylight "github.com/remshams/device-control/tui/pages/keylight"
	keylight_details "github.com/remshams/device-control/tui/pages/keylight/details"
	keylight_edit "github.com/remshams/device-control/tui/pages/keylight/edit"
	keylight_list "github.com/remshams/device-control/tui/pages/keylight/list"
	"github.com/remshams/device-control/keylight-control/control"

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

type Model struct {
	keylightAdapter *keylight.KeylightAdapter
	state           viewState
	keylights       []control.Keylight
	hueGroups       []hue_control.Group
	list            keylight_list.Model
	details         *keylight_details.Model
	edit            *keylight_edit.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{
		keylightAdapter: keylightAdapter,
		state:           initial,
		details:         nil,
		edit:            nil,
	}
}

func (m Model) Init() tea.Cmd {
	return m.init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		m.keylights = m.keylightAdapter.Control.Keylights()
		m.list = keylight_list.InitModel(m.keylightAdapter, m.keylights)
		m.state = list
	case pages_keylight.ReloadKeylights:
		m.keylightAdapter.Control.LoadOrDiscoverKeylights()
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
			cmd = tea.Batch(toast.CreateInfoToastAction("Keylight deleted"), pages_keylight.CreateReloadKeylights())
		}
	case keylight_list.ReloadKeylights:
		cmd = pages_keylight.CreateReloadKeylights()
	case pages_keylight.BackToListAction:
		m.details = nil
		m.edit = nil
		m.state = list
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
	switch m.state {
	case initial:
		return "Loading..."
	case list:
		return m.list.View()
	case details:
		return m.details.View()
	case add:
		return m.edit.View()
	case edit:
		return m.edit.View()
	default:
		return "Error"
	}

}

func (m *Model) init() tea.Cmd {
	return func() tea.Msg {
		m.keylightAdapter.Control.LoadOrDiscoverKeylights()
		return initMsg{}
	}
}
