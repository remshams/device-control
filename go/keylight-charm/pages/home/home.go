package home

import (
	"keylight-charm/lights/keylight"
	keylight_home "keylight-charm/pages/keylight/home"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	keylight keylight_home.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{
		keylight: keylight_home.InitModel(keylightAdapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.keylight.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.keylight, cmd = m.keylight.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.keylight.View()
}
