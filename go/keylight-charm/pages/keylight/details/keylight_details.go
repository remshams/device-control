package keylight_details

import (
	"fmt"
	"keylight-charm/components/actions"
	"keylight-charm/keylight"
	keylight_content "keylight-charm/pages/keylight/details/content"
	keylight_footer "keylight-charm/pages/keylight/details/footer"
	keylight_header "keylight-charm/pages/keylight/details/header"
	keylight_model "keylight-charm/pages/keylight/details/model"
	"keylight-charm/styles"

	"keylight-control/control"

	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	header  keylight_header.Model
	content keylight_content.Model
	footer  keylight_footer.Model
	state   keylight_model.ViewState
}

func InitModel(keylight *control.Keylight, keylightAdapter *keylight.KeylightAdapter) Model {
	return Model{
		header:  keylight_header.InitModel(keylight),
		content: keylight_content.InitModel(keylight, keylightAdapter),
		footer:  keylight_footer.InitModel(),
		state:   keylight_model.Navigate,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := []tea.Cmd{}
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
				cmds = append(cmds, keylight_model.CreateUpdateKeylight())
			} else {
				cmds = append(cmds, actions.CreateAbortAction())
			}
		default:
			var contentCommand tea.Cmd
			m.content, contentCommand = m.content.Update(msg, m.state)
			cmds = append(cmds, contentCommand)
		}
	default:
		var contentCommand tea.Cmd
		var footerCommand tea.Cmd
		m.content, contentCommand = m.content.Update(msg, m.state)
		m.footer, footerCommand = m.footer.Update(msg)
		cmds = append(cmds, contentCommand, footerCommand)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	style := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf("%s\n%s\n%s", style.Render(m.header.View()), style.Render(m.content.View(m.state)), m.footer.View(m.state))
}
