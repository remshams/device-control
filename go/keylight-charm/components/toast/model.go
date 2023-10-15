package toast

import (
	"keylight-charm/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type Toast interface {
	Message() string
}

type WarningToast struct {
	message string
}

func InitWarningToast(message string) WarningToast {
	return WarningToast{message}
}

func (warningToast WarningToast) Message() string {
	return styles.TextWarningColor.Render(warningToast.message)
}

type ErrorToast struct {
	message string
}

func InitErrorToast(message string) ErrorToast {
	return ErrorToast{message}
}

func (errorToast ErrorToast) Message() string {
	return styles.TextErrorColor.Render(errorToast.message)
}

func CreateWarningToastAction(message string) tea.Cmd {
	return func() tea.Msg {
		return WarningToast{message}
	}
}

func CreateErrorToastAction(message string) tea.Cmd {
	return func() tea.Msg {
		return ErrorToast{message}
	}
}
