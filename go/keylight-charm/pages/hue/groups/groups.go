package hue_groups

import (
	hue_control "hue-control/pubilc"

	tea "github.com/charmbracelet/bubbletea"
)

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

type GroupReloadedAction struct {
	Group hue_control.Group
}

func CreateGroupReloadedAction(group hue_control.Group) tea.Cmd {
	return func() tea.Msg {
		return GroupReloadedAction{
			Group: group,
		}
	}
}

type GroupsReloadedAction struct {
	Groups []hue_control.Group
}

func CreateGroupsReloadedAction(groups []hue_control.Group) tea.Cmd {
	return func() tea.Msg {
		return GroupsReloadedAction{
			Groups: groups,
		}
	}
}
