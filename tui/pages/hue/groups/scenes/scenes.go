package hue_group_scenes

import (
	hue_control "hue-control/pubilc"
	"github.com/remshams/device-control/tui/lights/hue"
	"github.com/remshams/device-control/tui/pages"
	hue_groups "github.com/remshams/device-control/tui/pages/hue/groups"
	"github.com/remshams/device-control/tui/stores"
	"github.com/remshams/device-control/tui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SceneSelectedAction struct {
	Scene hue_control.Scene
}

func CreateSceneSelectedAction(scene hue_control.Scene) tea.Cmd {
	return func() tea.Msg {
		return SceneSelectedAction{
			Scene: scene,
		}
	}
}

type sceneItem struct {
	title string
}

func (sceneItem sceneItem) Title() string {
	return sceneItem.title
}

func (sceneItem sceneItem) Description() string {
	return ""
}

func (sceneItem sceneItem) FilterValue() string {
	return sceneItem.title
}

type Model struct {
	adapter *hue.HueAdapter
	group   *hue_control.Group
	scenes  list.Model
}

func InitModel(adapter *hue.HueAdapter, group *hue_control.Group) Model {
	model := Model{
		adapter: adapter,
		group:   group,
		scenes:  createScenes(group.GetScenes()),
	}
	updateScenesLayout(&model.scenes)
	return model
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case pages.WindowResizeAction:
		updateScenesLayout(&m.scenes)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// This avoids conflicts between the keybindings of the main page
			// and the keybindings of the bubbletea list.
			// It ignores the enter and esc keybindings when the list is filtering or
			// in a filtered state.
			if m.scenes.FilterState() == list.Unfiltered {
				cmd = hue_groups.CreateBackToGroupDetailsAction()
			} else {
				m.scenes, cmd = m.scenes.Update(msg)
			}
		case "enter":
			if m.scenes.FilterState() == list.Filtering {
				m.scenes, cmd = m.scenes.Update(msg)
			} else {
				scene := m.group.GetSceneByName(m.scenes.SelectedItem().FilterValue())
				cmd = CreateSceneSelectedAction(*scene)
			}
		default:
			m.scenes, cmd = m.scenes.Update(msg)
		}
	default:
		m.scenes, cmd = m.scenes.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	return styles.ListStyles.Render(m.scenes.View())
}

func createScenes(scenes []hue_control.Scene) list.Model {
	var items []list.Item
	for _, scene := range scenes {
		items = append(items, sceneItem{title: scene.Name()})
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Scenes"
	return list
}

func updateScenesLayout(scenes *list.Model) {
	h, v := styles.ListStyles.GetFrameSize()
	scenes.SetSize(stores.LayoutStore.Width-h, stores.LayoutStore.Height-v)
}

func (m Model) findScene(index int) hue_control.Scene {
	return m.group.GetScenes()[index]
}
