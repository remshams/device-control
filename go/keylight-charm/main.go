package main

import (
	"fmt"
	checkbox "keylight-charm/components"
	"keylight-charm/keylight"
	"keylight-control/control"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type model struct {
	on      checkbox.Model
	cursor  int
	control *control.KeylightControl
}

func initModel(control control.KeylightControl) model {
	keylight := control.KeylightWithId(0)
	if keylight == nil {
		log.Error().Msg("No keylight found")
		os.Exit(1)
	}
	model := model{on: checkbox.New("On: ", keylight.Light.On), cursor: 0, control: &control}
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
	case control.KeylightCommand:
		m.control.SendKeylightCommand(msg)
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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	control := keylight.InitKeylightControl()
	p := tea.NewProgram(initModel(control))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
