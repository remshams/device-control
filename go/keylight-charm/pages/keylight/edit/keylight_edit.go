package keylight_edit

import (
	"fmt"
	keylight_model "keylight-charm/pages/keylight/details/model"
	"keylight-charm/styles"
	"keylight-control/control"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state string

const (
	navigate state = "view"
	insert         = "edit"
)

type Model struct {
	keylight *control.Keylight
	name     textinput.Model
	ip       textinput.Model
	port     textinput.Model
	cursor   int
	state    state
}

func InitModel(keylight *control.Keylight) Model {
	name := textinput.New()
	ip := textinput.New()
	port := textinput.New()
	if keylight != nil {
		name.SetValue(keylight.Metadata.Name)
		ip.SetValue(fmt.Sprintf("%s", keylight.Metadata.Ip))
		port.SetValue(strconv.Itoa(keylight.Metadata.Port))
	}
	cursor := 0
	state := navigate
	model := Model{keylight, name, ip, port, cursor, state}
	model.updateChildren()
	return model
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == navigate {
			cmd = m.processNavigateUpdate(msg)
		} else {
			cmd = m.processInsertUpdate(msg)
		}
	}

	return m, cmd
}

func (m *Model) processNavigateUpdate(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "j":
		m.increaseCursor()
	case "k":
		m.decreaseCursor()
	case "i":
		m.state = insert
	case "esc":
		cmd = keylight_model.CreateAbortAction()
	}
	m.updateChildren()
	return cmd
}

func (m *Model) processInsertUpdate(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "esc":
		m.state = navigate
	default:
		switch m.cursor {
		case 0:
			name, nameCmd := m.name.Update(msg)
			m.name = name
			cmd = nameCmd
		case 1:
			ip, ipCmd := m.ip.Update(msg)
			m.ip = ip
			cmd = ipCmd
		case 2:
			port, portCmd := m.port.Update(msg)
			m.port = port
			cmd = portCmd
		}

	}
	return cmd
}

func (m *Model) increaseCursor() {
	if m.cursor == 2 {
		m.cursor = 0
	} else {
		m.cursor++
	}
}

func (m *Model) decreaseCursor() {
	if m.cursor == 0 {
		m.cursor = 2
	} else {
		m.cursor--
	}
}

func (m *Model) updateChildren() {
	m.name.Blur()
	m.ip.Blur()
	m.port.Blur()
	switch m.cursor {
	case 0:
		m.name.Focus()
	case 1:
		m.ip.Focus()
	case 2:
		m.port.Focus()
	}
}

func (m Model) View() string {
	style := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf("%s\n%s\n%s\n%s", style.Render("Create/Edit keylight"), style.Render(m.renderLine("Name", m.name.View())), style.Render(m.renderLine("Ip", m.ip.View())), m.renderLine("Port", m.port.View()))
}

func (m *Model) renderLine(label string, value string) string {
	return fmt.Sprintf("%s %s", label, value)
}
