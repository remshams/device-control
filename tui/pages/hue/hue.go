package pages_hue

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ReloadBridgesAction struct{}

func CreateReloadBridgesAction() tea.Cmd {
	return func() tea.Msg {
		return ReloadBridgesAction{}
	}
}

type BridgesReloadedAction struct{}

func CreateBridgesReloadedAction() tea.Cmd {
	return func() tea.Msg {
		return BridgesReloadedAction{}
	}
}

type BackToHueHomeAction struct{}

func CreateBackToHueHomeAction() tea.Cmd {
	return func() tea.Msg {
		return BackToHueHomeAction{}
	}
}
