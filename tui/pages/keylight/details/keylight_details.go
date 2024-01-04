package keylight_details

import (
	"fmt"
	"github.com/remshams/device-control/tui/lights/keylight"
	pages_keylight "github.com/remshams/device-control/tui/pages/keylight"
	keylight_content "github.com/remshams/device-control/tui/pages/keylight/details/content"
	keylight_footer "github.com/remshams/device-control/tui/pages/keylight/details/footer"
	keylight_header "github.com/remshams/device-control/tui/pages/keylight/details/header"
	keylight_model "github.com/remshams/device-control/tui/pages/keylight/details/model"
	"github.com/remshams/device-control/tui/styles"

	"github.com/remshams/device-control/keylight-control/control"

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
	keylight.LoadLights()
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
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			m.state = keylight_model.Insert
		case "enter":
			cmds = append(cmds, m.updateContent(msg))
			m.state = keylight_model.Navigate
		case "esc":
			if m.state == keylight_model.Insert {
				cmds = append(cmds, m.updateContent(msg))
				m.state = keylight_model.Navigate
			} else {
				cmds = append(cmds, pages_keylight.CreateBackToListAction())
			}
		default:
			cmds = append(cmds, m.updateContent(msg))
		}
	default:
		var footerCommand tea.Cmd
		m.footer, footerCommand = m.footer.Update(msg)
		cmds = append(cmds, m.updateContent(msg), footerCommand)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) updateContent(msg tea.Msg) tea.Cmd {
	var contentCommand tea.Cmd
	m.content, contentCommand = m.content.Update(msg, m.state)
	return contentCommand
}

func (m Model) View() string {
	style := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf("%s\n%s\n%s", style.Render(m.header.View()), style.Render(m.content.View(m.state)), m.footer.View(m.state))
}
