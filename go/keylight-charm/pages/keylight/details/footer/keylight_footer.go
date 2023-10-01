package keylight_footer

import (
	"fmt"
	keylight_model "keylight-charm/pages/keylight/details/model"
	"keylight-charm/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	statusMessage string
}

func InitModel() Model {
	return Model{statusMessage: ""}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case keylight_model.CommandResult:
		m.statusMessage = m.createStatusMessage(msg.Status)
	}
	return m, nil
}

func (m Model) View(state keylight_model.ViewState) string {
	style := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf("Mode: %s\nStatus: %s", style.Render(fmt.Sprintf("%s", state)), m.statusMessage)
}

func (m *Model) createStatusMessage(status keylight_model.CommandStatus) string {
	switch status {
	case keylight_model.Success:
		return "Light values set"
	case keylight_model.Error:
		return "Could not set light values"
	default:
		return ""
	}
}
