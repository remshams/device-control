package main

import (
	"fmt"
	checkbox "keylight-charm/components"

	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	on     checkbox.Model
	cursor int
}

func initModel() model {
	model := model{on: checkbox.New("On: "), cursor: 0}
	model.selectedElement()
	return model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "j", "down":
			m.cursor++
			m.selectedElement()
		case "k", "up":
			m.cursor--
			m.selectedElement()
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			cmd = m.updateChild(msg)
		}
	}
	return m, cmd
}

func (m *model) updateChild(msg tea.Msg) tea.Cmd {
	switch m.cursor {
	case 0:
		on, cmd := m.on.Update(msg)
		m.on = on
		return cmd
	default:
		return nil
	}
}

func (m *model) selectedElement() {
	switch m.cursor {
	case 0:
		m.on.Focus = true
	default:
		m.on.Focus = false
	}
}

func (m model) View() string {
	title := "Update keylight"
	cursor := " "
	if m.on.Focus {
		cursor = ">"
	}
	on := fmt.Sprintf("%s %s", cursor, m.on.View())

	return fmt.Sprintf("%s \n\n %s", title, on)

}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
