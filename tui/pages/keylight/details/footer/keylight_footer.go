package keylight_footer

import (
	"fmt"
	keylight_model "github.com/remshams/device-control/tui/pages/keylight/details/model"
	"github.com/remshams/device-control/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
}

func InitModel() Model {
	return Model{}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View(state keylight_model.ViewState) string {
	style := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf("Mode: %s", style.Render(fmt.Sprintf("%s", state)))
}
