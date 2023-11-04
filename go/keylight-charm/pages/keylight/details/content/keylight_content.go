package keylight_content

import (
	"fmt"
	"keylight-charm/components/checkbox"
	"keylight-charm/components/textinput"
	"keylight-charm/components/toast"
	"keylight-charm/keylight"
	"keylight-charm/pages/keylight/details/model"
	"keylight-charm/styles"
	"keylight-control/control"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type Model struct {
	keylight        *control.Keylight
	on              checkbox.Model
	brightness      textinput.Model
	temperature     textinput.Model
	cursor          int
	keylightAdapter *keylight.KeylightAdapter
}

func InitModel(keylight *control.Keylight, keylightAdapter *keylight.KeylightAdapter) Model {
	model := Model{
		keylight:    keylight,
		on:          checkbox.New("On: ", false),
		brightness:  kl_textinput.CreateTextInputModel(),
		temperature: kl_textinput.CreateTextInputModel(),
		cursor:      0, keylightAdapter: keylightAdapter,
	}
	model.updateKeylight()
	return model
}

func (m Model) Update(msg tea.Msg, state keylight_model.ViewState) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if state == keylight_model.Insert {
			cmd = m.processInInsertMode(msg)
		} else {
			cmd = m.processInNavigateMode(msg)
		}
	}
	return m, cmd
}

func (m *Model) processInInsertMode(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "enter":
		cmd = m.sendCommand()
	case "esc":
		m.updateKeylight()
		fmt.Println("stop polling")
	default:
		cmd = m.updateChild(msg)
	}
	return cmd
}

func (m *Model) processInNavigateMode(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "j", "down":
		m.increaseCursor()
		m.selectedElement()
	case "k", "up":
		m.decreaseCursor()
		m.selectedElement()
	case "enter":
		cmd = m.sendCommand()
	}
	return cmd
}

func (m *Model) updateChild(msg tea.Msg) tea.Cmd {
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

func (m *Model) selectedElement() {
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

func (m *Model) increaseCursor() {
	m.cursor++
	if m.cursor > 2 {
		m.cursor = 0
	}
}

func (m *Model) decreaseCursor() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = 2
	}
}

func (m Model) View(state keylight_model.ViewState) string {
	on := fmt.Sprintf("%s", m.on.View())
	brightness := kl_textinput.CreateTextInputView(m.brightness, "Brightness", "%")
	temperature := kl_textinput.CreateTextInputView(m.temperature, "Temperature", "")
	on = m.renderLine(on, m.cursor == 0, state == keylight_model.Insert)
	brightness = m.renderLine(brightness, m.cursor == 1, state == keylight_model.Insert)
	temperature = m.renderLine(temperature, m.cursor == 2, state == keylight_model.Insert)

	return fmt.Sprintf("%s\n\n%s\n\n%s", on, brightness, temperature)
}

func (m *Model) renderLine(line string, isActive bool, isEdit bool) string {
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

func (m *Model) sendCommand() tea.Cmd {
	err := m.keylightAdapter.SendCommand(m.keylight.Metadata.Id, m.on.Checked, m.brightness.Value(), m.temperature.Value())
	var status keylight_model.CommandStatus
	if err != nil {
		status = keylight_model.Error
	} else {
		status = keylight_model.Success
	}
	return toast.CreateInfoToastAction(m.createStatusMessage(status))
}

func (m *Model) updateKeylight() {
	keylight := m.keylightAdapter.KeylightControl.KeylightWithId(0)
	keylight.LoadLights()
	if keylight == nil {
		log.Error().Msg("No keylight found")
		os.Exit(1)
	}
	m.on = checkbox.New("On: ", keylight.Light.On)
	m.brightness.SetValue(fmt.Sprintf("%d", keylight.Light.Brightness))
	m.temperature.SetValue(fmt.Sprintf("%d", keylight.Light.Temperature))
	m.selectedElement()
}

func (m Model) createStatusMessage(commandStatus keylight_model.CommandStatus) string {
	switch commandStatus {
	case keylight_model.Success:
		return "Light values set"
	case keylight_model.Error:
		return "Could not set light values"
	default:
		return ""
	}

}
