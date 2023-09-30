package keylight_details

import (
	"fmt"
	"keylight-charm/keylight"
	keylight_content "keylight-charm/pages/keylight/details/content"
	keylight_header "keylight-charm/pages/keylight/details/header"
	"keylight-charm/styles"

	"keylight-control/control"

	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	header  keylight_header.Model
	content keylight_content.Model
}

func InitModel(keylight *control.Keylight, keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{
		header:  keylight_header.InitModel(keylight),
		content: keylight_content.InitModel(keylight, keylightAdapter),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.content, cmd = m.content.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	style := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf("%s\n%s", style.Render(m.header.View()), m.content.View())
}
