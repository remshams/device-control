package hue_groups

import tea "github.com/charmbracelet/bubbletea"

type BackToListAction struct{}

func CreateBackToListAction() tea.Cmd {
	return func() tea.Msg {
		return BackToListAction{}
	}
}

type BackToGroupDetailsAction struct{}

func CreateBackToGroupDetailsAction() tea.Cmd {
	return func() tea.Msg {
		return BackToGroupDetailsAction{}
	}
}

type BackToGroupHomeAction struct{}

func CreateBackToGroupHomeAction() tea.Cmd {
	return func() tea.Msg {
		return BackToGroupHomeAction{}
	}
}
