package page_help

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	help   help.Model
	keyMap *help.KeyMap
}

func New() Model {
	return Model{
		help:   help.New(),
		keyMap: nil,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetKeyMapMsg:
		m.keyMap = &msg.keyMap
	case ResetKeyMapMsg:
		m.keyMap = nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.keyMap == nil {
		return ""
	}
	return m.help.View(*m.keyMap)
}
