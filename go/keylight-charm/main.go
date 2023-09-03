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

type initMsg struct{}

type model struct {
	isInit       bool
	on           checkbox.Model
	brightness   textinput.Model
	cursor       int
	isInsertMode bool
	control      *control.KeylightControl
}

func initModel(control control.KeylightControl) model {
	model := model{isInit: false, on: checkbox.New("On: ", false), brightness: textinput.New(), cursor: 0, isInsertMode: false, control: &control}
	return model
}

func (m model) Init() tea.Cmd {
	return m.discoverKeylights()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		m.processInit()
	case tea.KeyMsg:
		if m.isInsertMode {
			cmd = m.processInInsertMode(msg)
		} else {
			cmd = m.processInNavigateMode(msg)
		}
	}
	return m, cmd
}

func (m *model) processInit() {
	keylight := m.control.KeylightWithId(0)
	if keylight == nil {
		log.Error().Msg("No keylight found")
		os.Exit(1)
	}
	brightness := textinput.New()
	brightness.SetValue(fmt.Sprintf("%d", keylight.Light.Brightness))
	m.brightness = brightness
	m.on = checkbox.New("On: ", keylight.Light.On)
	m.isInit = true
	m.selectedElement()
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
		m.isInsertMode = true
	case "j", "down":
		m.increaseCursor()
		m.selectedElement()
	case "k", "up":
		m.decreaseCursor()
		m.selectedElement()
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

func (m *model) increaseCursor() {
	m.cursor++
	if m.cursor > 1 {
		m.cursor = 0
	}
}

func (m *model) decreaseCursor() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = 1
	}
}

func (m model) View() string {
	title := "Update keylight"
	if m.isInit {
		on := fmt.Sprintf("%s", m.on.View())
		brightness := fmt.Sprintf("Brightness %s%%", m.brightness.View())
		lines := m.renderCursor([]string{on, brightness})

		return fmt.Sprintf("%s \n\n %s \n\n %s", title, lines[0], lines[1])
	} else {
		return fmt.Sprintf("%s \n\n %s", title, "Loading...")
	}
}

func (m *model) renderCursor(lines []string) []string {
	var linesWithSelector []string
	var editMarker string
	if m.isInsertMode {
		editMarker = "(edit)"
	} else {
		editMarker = ""
	}
	for i, line := range lines {
		var newLine string
		if i == m.cursor {
			newLine = fmt.Sprintf("> %s %s", line, editMarker)
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

func (m *model) discoverKeylights() tea.Cmd {
	return func() tea.Msg {
		m.control.LoadOrDiscoverKeylights()
		return initMsg{}
	}
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
