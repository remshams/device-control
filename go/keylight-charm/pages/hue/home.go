package hue_home

import tea "github.com/charmbracelet/bubbletea"

type Model struct{}

func InitModel() Model {
	return Model{}
}

func (m Model) Update(mgs tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return "HueLights"
}
