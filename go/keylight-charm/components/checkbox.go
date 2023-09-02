package checkbox

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Focus   bool
	Label   string
	Checked bool
}

func New(label string, checked bool) Model {
	return Model{Focus: false, Label: label, Checked: checked}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.Focus {
		switch mgs := msg.(type) {
		case tea.KeyMsg:
			switch mgs.String() {
			case " ":
				m.Checked = !m.Checked
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	checked := " "
	if m.Checked {
		checked = "x"
	}

	return fmt.Sprintf("%s [%s]", m.Label, checked)
}
