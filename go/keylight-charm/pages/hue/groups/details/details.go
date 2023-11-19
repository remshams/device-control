package hue_group_details

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/components/checkbox"
	kl_cursor "keylight-charm/components/cursor"
	"keylight-charm/lights/hue"
	hue_groups "keylight-charm/pages/hue/groups"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState = string

const (
	navigate viewState = "navigate"
	insert   viewState = "insert"
)

type Model struct {
	adapter *hue.HueAdapter
	group   hue_control.Group
	on      checkbox.Model
	state   viewState
	cursor  int
}

func InitModel(adapter *hue.HueAdapter, group hue_control.Group) Model {
	return Model{
		adapter: adapter,
		group:   group,
		on:      checkbox.New("On", group.GetOn()),
		state:   navigate,
		cursor:  0,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			m.state = insert
			m.focuesView()
		case "enter":
			m.state = navigate
			m.unfocusView()
			m.sendGroup()
		case "esc":
			if m.state == navigate {
				cmd = hue_groups.CreateBackToListAction()
			} else {
				m.state = navigate
				m.unfocusView()
				m.resetView()
			}
		default:
			m.on, cmd = m.on.Update(msg)
		}

	}
	return m, cmd
}

func (m Model) View() string {
	return kl_cursor.RenderLine(m.on.View(), true, m.state == insert)
}

func (m *Model) focuesView() {
	m.on.Focus = true
}

func (m *Model) unfocusView() {
	m.on.Focus = false
}

func (m *Model) resetView() {
	m.on.Checked = m.group.GetOn()
}

func (m *Model) updateGroup() {
	m.group.SetOn(m.on.Checked)
}

func (m *Model) sendGroup() {
	m.updateGroup()
	m.group.SendGroup()
}
