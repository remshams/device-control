package toast

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	toast Toast
}

func InitModel() Model {
	return Model{}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrorToast:
		m.toast = msg
	case WarningToast:
		m.toast = msg
	default:
		m.toast = nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.toast != nil {
		return m.toast.Message()
	} else {
		return ""
	}
}

func RenderToast(toast Toast) string {
	return toast.Message()
}