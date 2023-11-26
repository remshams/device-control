package hue_group_scenes

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"
	hue_groups "keylight-charm/pages/hue/groups"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	adapter *hue.HueAdapter
	group   hue_control.Group
	menu    list.Model
}

func InitModel(adapter *hue.HueAdapter, group hue_control.Group) Model {
	return Model{
		adapter: adapter,
		group:   group,
		menu:    list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = hue_groups.CreateBackToGroupDetailsAction()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return "List of scenes"
}
