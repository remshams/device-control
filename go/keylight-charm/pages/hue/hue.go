package pages_hue

import (
	hue_control "hue-control/pubilc"

	tea "github.com/charmbracelet/bubbletea"
)

type ReloadBridgesAction struct{}

func CreateReloadBridgesAction() tea.Cmd {
	return func() tea.Msg {
		return ReloadBridgesAction{}
	}
}

type BridgesReloadedAction struct {
	Bridges []hue_control.Bridge
}

func CreateBridgesReloadedAction(bridges []hue_control.Bridge) tea.Cmd {
	return func() tea.Msg {
		return BridgesReloadedAction{
			Bridges: bridges,
		}
	}
}

type BackToHueHomeAction struct{}

func CreateBackToHueHomeAction() tea.Cmd {
	return func() tea.Msg {
		return BackToHueHomeAction{}
	}
}
