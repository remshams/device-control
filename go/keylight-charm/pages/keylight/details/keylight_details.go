package keylight_details

import (
	"fmt"
	"keylight-charm/keylight"
	keylight_content "keylight-charm/pages/keylight/details/content"
	keylight_header "keylight-charm/pages/keylight/details/header"
	keylight_model "keylight-charm/pages/keylight/details/model"
	"keylight-charm/styles"

	"keylight-control/control"

	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type AbortAction struct{}

type Model struct {
	header  keylight_header.Model
	content keylight_content.Model
	state   keylight_model.ViewState
}

func InitModel(keylight *control.Keylight, keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{
		header:  keylight_header.InitModel(keylight),
		content: keylight_content.InitModel(keylight, keylightAdapter),
		state:   keylight_model.Navigate,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case keylight_content.StateChanged:
		m.state = msg.State
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			m.state = keylight_model.Insert
		case "esc":
			if m.state == keylight_model.Insert {
				m.state = keylight_model.Navigate
				cmd = keylight_model.CreateUpdateKeylight()
			} else {
				cmd = m.abortAction()
			}
		default:
			m.content, cmd = m.content.Update(msg, m.state)
		}
	default:
		m.content, cmd = m.content.Update(msg, m.state)
	}
	return m, cmd
}

func (m Model) View() string {
	style := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf("%s\n%s", style.Render(m.header.View()), m.content.View(m.state))
}

func (m *Model) abortAction() tea.Cmd {
	return func() tea.Msg {
		return AbortAction{}
	}
}
