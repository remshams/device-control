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

func (warningToast WarningToast) Message() string {
	return styles.TextWarningColor.Render(warningToast.message)
}

func CreateWarningToastAction(message string) tea.Cmd {
	return func() tea.Msg {
		return WarningToast{message}
	}
}

type ErrorToast struct {
	message string
}

func (errorToast ErrorToast) Message() string {
	return styles.TextErrorColor.Render(errorToast.message)
}

func CreateErrorToastAction(message string) tea.Cmd {
	return func() tea.Msg {
		return ErrorToast{message}
	}
}

type InfoToast struct {
	message string
}

func (infoToast InfoToast) Message() string {
	return styles.TextInfoColor.Render(infoToast.message)
}

func CreateInfoToastAction(message string) tea.Cmd {
	return func() tea.Msg {
		return InfoToast{message}
	}
}
