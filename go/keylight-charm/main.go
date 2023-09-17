package main

import (
	"fmt"
	"keylight-charm/components/checkbox"
	"keylight-charm/components/textinput"
	"keylight-charm/keylight"
	"keylight-charm/styles"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type viewState string

const (
	initial  viewState = "initial"
	edit               = "edit"
	navigate           = "navigate"
	inError            = "error"
)

type initMsg struct{}

type model struct {
	state           viewState
	on              checkbox.Model
	brightness      textinput.Model
	temperature     textinput.Model
	cursor          int
	keylightAdapter *keylight.KeylightAdapter
	message         string
}

func initModel(keylightAdapter keylight.KeylightAdapter) model {
	model := model{state: initial, on: checkbox.New("On: ", false), brightness: kl_textinput.CreateTextInputModel(), temperature: kl_textinput.CreateTextInputModel(), cursor: 0, keylightAdapter: &keylightAdapter, message: ""}
	return model
}

func (m model) Init() tea.Cmd {
	return m.discoverKeylights()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initMsg:
		m.updateKeylight()
	case tea.KeyMsg:
		if m.state == edit {
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
		m.state = navigate
		m.updateKeylight()
	case "enter":
		m.state = navigate
		m.sendCommand()
	default:
		cmd = m.updateChild(msg)
	}
	return cmd
}

func (m *model) processInNavigateMode(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "i":
		m.state = edit
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
	case 2:
		temperature, cmd := m.temperature.Update(msg)
		m.temperature = temperature
		return cmd
	default:
		return nil
	}
}

func (m *model) selectedElement() {
	m.on.Focus = false
	m.brightness.Blur()
	m.temperature.Blur()
	switch m.cursor {
	case 0:
		m.on.Focus = true
	case 1:
		m.brightness.Focus()
	case 2:
		m.temperature.Focus()
	}
}

func (m *model) increaseCursor() {
	m.cursor++
	if m.cursor > 2 {
		m.cursor = 0
	}
}

func (m *model) decreaseCursor() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = 2
	}
}

func (m model) View() string {
	title := "Update keylight"
	if m.state != initial {
		on := fmt.Sprintf("%s", m.on.View())
		brightness := kl_textinput.CreateTextInputView(m.brightness, "Brightness", "%")
		temperature := kl_textinput.CreateTextInputView(m.temperature, "Temperature", "")
		on = m.renderLine(on, m.cursor == 0, m.state == edit)
		brightness = m.renderLine(brightness, m.cursor == 1, m.state == edit)
		temperature = m.renderLine(temperature, m.cursor == 2, m.state == edit)

		return fmt.Sprintf("%s \n\n %s \n\n %s \n\n %s \n\n\n Mode: %s \n\n\n Status: %s", title, on, brightness, temperature, m.state, m.message)
	} else {
		return fmt.Sprintf("%s \n\n %s", title, "Loading...")
	}
}

func (m *model) renderLine(line string, isActive bool, isEdit bool) string {
	style := lipgloss.NewStyle().PaddingLeft(styles.Padding)
	cursor := ""
	if isActive {
		style = style.UnsetPaddingLeft()
		cursorStyles := lipgloss.NewStyle()
		cursorStyles.Foreground(styles.AccentColor)
		cursor = styles.TextAccentColor.Render(">")
	}
	edit := ""
	if isActive && isEdit {
		edit = "(edit)"
	}
	return style.Render(fmt.Sprintf("%s %s %s", cursor, line, edit))
}

func (m *model) sendCommand() {
	err := m.keylightAdapter.SendCommand(m.on.Checked, m.brightness.Value(), m.temperature.Value())
	if err != nil {
		m.message = "Could not set light values"
	} else {
		m.message = "Light values set"
	}
	m.updateKeylight()
}

func (m *model) discoverKeylights() tea.Cmd {
	return func() tea.Msg {
		m.keylightAdapter.Control.LoadOrDiscoverKeylights()
		return initMsg{}
	}
}

func (m *model) updateKeylight() {
	keylight := m.keylightAdapter.Control.KeylightWithId(0)
	if keylight == nil {
		log.Error().Msg("No keylight found")
		os.Exit(1)
	}
	m.on = checkbox.New("On: ", keylight.Light.On)
	m.brightness.SetValue(fmt.Sprintf("%d", keylight.Light.Brightness))
	m.temperature.SetValue(fmt.Sprintf("%d", keylight.Light.Temperature))
	m.state = navigate
	m.selectedElement()
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	keylightAdapter := keylight.NewKeylightAdapter()
	p := tea.NewProgram(initModel(keylightAdapter))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
