package main

import (
	"fmt"
	checkbox "keylight-charm/components"
	"keylight-charm/keylight"
	"keylight-control/control"
	"strconv"

	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type model struct {
	on         checkbox.Model
	brightness textinput.Model
	cursor     int
	control    *control.KeylightControl
}

func initModel(control control.KeylightControl) model {
	keylight := control.KeylightWithId(0)
	if keylight == nil {
		log.Error().Msg("No keylight found")
		os.Exit(1)
	}
	brightness := textinput.New()
	brightness.SetValue(fmt.Sprintf("%d", keylight.Light.Brightness))
	model := model{on: checkbox.New("On: ", keylight.Light.On), brightness: brightness, cursor: 0, control: &control}
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

		cmd = m.updateChild(msg)
		switch msg.String() {
		case "j", "down":
			m.cursor++
			m.selectedElement()
		case "k", "up":
			m.cursor--
			m.selectedElement()
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.sendCommand()
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
	case 1:
		brightness, cmd := m.brightness.Update(msg)
		m.brightness = brightness
		return cmd
	default:
		return nil
	}
}

func (m *model) selectedElement() {
	switch m.cursor {
	case 0:
		m.on.Focus = true
		m.brightness.Blur()
	case 1:
		m.brightness.Focus()
		m.on.Focus = false
	default:
		m.on.Focus = false
		m.brightness.Blur()
	}
}

func (m model) View() string {
	title := "Update keylight"
	onCursor := " "
	if m.on.Focus {
		onCursor = ">"
	}
	brightnessCursor := " "
	if m.brightness.Focused() {
		brightnessCursor = ">"
	}
	on := fmt.Sprintf("%s %s", onCursor, m.on.View())
	brightness := fmt.Sprintf("%s Brightness %s%%", brightnessCursor, m.brightness.View())

	return fmt.Sprintf("%s \n\n %s \n\n %s", title, on, brightness)

}

func (m *model) sendCommand() {
	on := m.on.Checked
	brightness, _ := strconv.Atoi(m.brightness.Value())
	m.control.SendKeylightCommand(control.KeylightCommand{Id: 0, Command: control.LightCommand{On: &on, Brightness: &brightness}})
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
