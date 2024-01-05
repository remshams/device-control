package page_title

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/device-control/tui/styles"
)

type Model struct {
	header string
}

func New() Model {
	return Model{header: ""}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case SetPageTitleMsg:
		m.header = msg
	}
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		Bold(true).
		Background(styles.HeadlineBackgroundColor).
		Render(m.header)
}
