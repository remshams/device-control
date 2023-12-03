package hue_group_scenes

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"
	"keylight-charm/pages"
	hue_groups "keylight-charm/pages/hue/groups"
	"keylight-charm/stores"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var scenesStyle = lipgloss.NewStyle().Margin(1, 2)

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
	group   hue_control.Group
	scenes  list.Model
}

func InitModel(adapter *hue.HueAdapter, group hue_control.Group) Model {
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
			cmd = hue_groups.CreateBackToGroupDetailsAction()
		case "enter":
			scene := m.findScene(m.scenes.Index())
			cmd = CreateSceneSelectedAction(scene)
		default:
			m.scenes, cmd = m.scenes.Update(msg)
		}
	default:
		m.scenes, cmd = m.scenes.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	return scenesStyle.Render(m.scenes.View())
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
	h, v := scenesStyle.GetFrameSize()
	scenes.SetSize(stores.LayoutStore.Width-h, stores.LayoutStore.Height-v)
}

func (m Model) findScene(index int) hue_control.Scene {
	return m.group.GetScenes()[index]
}
