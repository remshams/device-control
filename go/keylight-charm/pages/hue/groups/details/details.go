package hue_group_details

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/components/checkbox"
	"keylight-charm/lights/hue"
	hue_groups "keylight-charm/pages/hue/groups"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	adapter *hue.HueAdapter
	group   hue_control.Group
	on      checkbox.Model
}

func InitModel(adapter *hue.HueAdapter, group hue_control.Group) Model {
	return Model{
		adapter: adapter,
		group:   group,
		on:      checkbox.New("On", group.GetOn()),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			cmd = hue_groups.CreateBackToListAction()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return m.on.View()
}
