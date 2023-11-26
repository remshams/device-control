package hue_group_details

import (
	"fmt"
	hue_control "hue-control/pubilc"
	"keylight-charm/components/checkbox"
	kl_cursor "keylight-charm/components/cursor"
	"keylight-charm/lights/hue"
	hue_groups "keylight-charm/pages/hue/groups"
	hue_group_scenes "keylight-charm/pages/hue/groups/scenes"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState = string

const (
	navigate viewState = "navigate"
	insert   viewState = "insert"
	scenes   viewState = "scenes"
)

type Model struct {
	adapter *hue.HueAdapter
	group   hue_control.Group
	on      checkbox.Model
	state   viewState
	cursor  int
	scenes  hue_group_scenes.Model
}

func InitModel(adapter *hue.HueAdapter, group hue_control.Group) Model {
	return Model{
		adapter: adapter,
		group:   group,
		on:      checkbox.New("On", group.GetOn()),
		state:   navigate,
		cursor:  0,
		scenes:  hue_group_scenes.InitModel(adapter, group),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case hue_groups.BackToGroupDetailsAction:
		m.state = navigate
	case tea.KeyMsg:
		switch m.state {
		case scenes:
			m.scenes, cmd = m.scenes.Update(msg)
		default:
			switch msg.String() {
			case "i":
				m.processInsert()
			case "k":
				m.incrementCursor()
			case "j":
				m.decrementCursor()
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

	}
	return m, cmd
}

func (m Model) View() string {
	if m.state == scenes {
		return m.scenes.View()
	} else {
		cursor := kl_cursor.RenderLine(m.on.View(), m.cursor == 0, m.state == insert)
		scenes := kl_cursor.RenderLine("Scenes", m.cursor == 1, m.state == insert)
		return fmt.Sprintf("%s\n\n%s", cursor, scenes)
	}
}

func (m *Model) processInsert() {
	if m.cursor == 1 {
		m.state = scenes
	} else {
		m.state = insert
		m.focusView()
	}
}

func (m *Model) incrementCursor() {
	m.cursor = (m.cursor + 1) % 2
}

func (m *Model) decrementCursor() {
	m.cursor = int(math.Abs(float64((m.cursor - 1) % 2)))
}

func (m *Model) focusView() {
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
