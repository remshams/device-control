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
	on           checkbox.Model
	brightness   textinput.Model
	cursor       int
	isInsertMode bool
	control      *control.KeylightControl
}

func initModel(control control.KeylightControl) model {
	keylight := control.KeylightWithId(0)
	if keylight == nil {
		log.Error().Msg("No keylight found")
		os.Exit(1)
	}
	brightness := textinput.New()
	brightness.SetValue(fmt.Sprintf("%d", keylight.Light.Brightness))
	model := model{on: checkbox.New("On: ", keylight.Light.On), brightness: brightness, cursor: 0, isInsertMode: false, control: &control}
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

		if m.isInsertMode {
			cmd = m.processInInsertMode(msg)
		} else {
			cmd = m.processInNavigateMode(msg)
		}
	}
	return m, cmd
}

func (m *model) processInInsertMode(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		m.isInsertMode = false
	default:
		cmd = m.updateChild(msg)
	}
	return cmd
}

func (m *model) processInNavigateMode(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "i":
		if !m.isInsertMode {
			m.isInsertMode = true
		}
	case "j", "down":
		if !m.isInsertMode {
			m.cursor++
			m.normalizeCursor()
			m.selectedElement()
		}
	case "k", "up":
		if !m.isInsertMode {
			m.cursor--
			m.normalizeCursor()
			m.selectedElement()
		}
	case "ctrl+c", "q":
		cmd = tea.Quit
	case "enter":
		m.sendCommand()
	}
	return cmd
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
	m.on.Focus = false
	m.brightness.Blur()
	switch m.cursor {
	case 0:
		m.on.Focus = true
	case 1:
		m.brightness.Focus()
	}
}

func (m *model) normalizeCursor() {
	if m.cursor < 0 {
		m.cursor = 1
	}
	if m.cursor > 1 {
		m.cursor = 0
	}
}

func (m model) View() string {
	title := "Update keylight"
	on := fmt.Sprintf("%s", m.on.View())
	brightness := fmt.Sprintf("Brightness %s%%", m.brightness.View())
	lines := m.renderCursor([]string{on, brightness})

	return fmt.Sprintf("%s \n\n %s \n\n %s", title, lines[0], lines[1])
}

func (m *model) renderCursor(lines []string) []string {
	var linesWithSelector []string
	for i, line := range lines {
		var newLine string
		if i == m.cursor {
			newLine = fmt.Sprintf("> %s", line)
		} else {
			newLine = fmt.Sprintf("  %s", line)
		}
		linesWithSelector = append(linesWithSelector, newLine)
	}
	return linesWithSelector
}

func (m *model) sendCommand() {
	on := m.on.Checked
	brightness, _ := strconv.Atoi(m.brightness.Value())
	m.control.SendKeylightCommand(control.KeylightCommand{Id: 0, Command: control.LightCommand{On: &on, Brightness: &brightness}})
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	control := keylight.InitKeylightControl()
	p := tea.NewProgram(initModel(control))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
