package hue_group_details

import (
	"fmt"
	hue_control "hue-control/pubilc"
	"ui/components/checkbox"
	kl_cursor "ui/components/cursor"
	"ui/components/toast"
	"ui/lights/hue"
	hue_groups "ui/pages/hue/groups"
	hue_group_scenes "ui/pages/hue/groups/scenes"
	"math"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState = string
type selectedSceneSent struct {
	success bool
	message string
}

const (
	navigate viewState = "navigate"
	scenes   viewState = "scenes"
)

type Model struct {
	adapter *hue.HueAdapter
	group   *hue_control.Group
	on      checkbox.Model
	state   viewState
	cursor  int
	scenes  hue_group_scenes.Model
}

func InitModel(adapter *hue.HueAdapter, group *hue_control.Group) Model {
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
		m.reloadLights()
	case hue_group_scenes.SceneSelectedAction:
		cmd = tea.Batch(toast.CreateInfoToastAction("Setting scene"), m.sendScene(msg.Scene))
	case selectedSceneSent:
		if msg.success {
			m.reloadLights()
			cmd = toast.CreateSuccessToastAction(msg.message)
		} else {
			cmd = toast.CreateErrorToastAction(msg.message)
		}
	case tea.KeyMsg:
		switch m.state {
		case scenes:
			m.scenes, cmd = m.scenes.Update(msg)
		default:
			switch msg.String() {
			case "k":
				m.incrementCursor()
			case "j":
				m.decrementCursor()
			case "enter":
				cmd = m.processEnterKey()
			case "esc":
				if m.state == navigate {
					cmd = hue_groups.CreateBackToListAction()
				} else {
					m.state = navigate
					m.resetView()
				}
			default:
				m.on, cmd = m.on.Update(msg)
			}
		}
	default:
		if m.state == scenes {
			m.scenes, cmd = m.scenes.Update(msg)
		}

	}
	return m, cmd
}

func (m *Model) processEnterKey() tea.Cmd {
	var cmd tea.Cmd
	switch m.cursor {
	case 0:
		m.group.SetOn(!m.group.GetOn())
		err := m.group.SendGroup()
		m.reloadLights()
		toastCmd := toast.CreateSuccessToastAction("Group updated")
		if err != nil {
			toastCmd = toast.CreateErrorToastAction("Error updating group")
		}
		cmd = toastCmd
	case 1:
		m.state = scenes
	}
	return cmd
}

func (m Model) View() string {
	if m.state == scenes {
		return m.scenes.View()
	} else {
		cursor := kl_cursor.RenderLine(m.on.View(), m.cursor == 0, false)
		scenes := kl_cursor.RenderLine("Scenes", m.cursor == 1, false)
		return fmt.Sprintf("%s\n\n%s", cursor, scenes)
	}
}

func (m *Model) incrementCursor() {
	m.cursor = (m.cursor + 1) % 2
}

func (m *Model) decrementCursor() {
	m.cursor = int(math.Abs(float64((m.cursor - 1) % 2)))
}

func (m *Model) sendScene(scene hue_control.Scene) tea.Cmd {
	err := m.group.SetScene(scene)
	return func() tea.Msg {
		if err == nil {
			return selectedSceneSent{
				success: true,
				message: "Scene set",
			}
		} else {
			return selectedSceneSent{
				success: false,
				message: "Error setting scene",
			}
		}
	}
}

func (m *Model) reloadLights() {
	m.group = m.adapter.Control.GetBridgeById(m.group.GetBridgeId()).GetGroupById(m.group.GetId())
	m.resetView()
}

func (m *Model) resetView() {
	m.on.Checked = m.group.GetOn()
}
