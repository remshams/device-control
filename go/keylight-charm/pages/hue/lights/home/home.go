package hue_lights_home

import (
	hue_control "hue-control/pubilc"
	"keylight-charm/lights/hue"
	hue_lights "keylight-charm/pages/hue/lights"
	hue_lights_details "keylight-charm/pages/hue/lights/details"
	hue_lights_list "keylight-charm/pages/hue/lights/list"

	tea "github.com/charmbracelet/bubbletea"
)

type viewState string

const (
	list    viewState = "list"
	details viewState = "details"
)

type Model struct {
	adapter *hue.HueAdapter
	state   viewState
	list    hue_lights_list.Model
	details *hue_lights_details.Model
}

func InitModel(adapter *hue.HueAdapter) Model {
	return Model{
		adapter: adapter,
		state:   list,
		list:    hue_lights_list.InitModel(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	switch msg := msg.(type) {
	case hue_lights.BackToLightHomeAction:
		m.details = nil
		m.state = list
	case hue_lights_list.LightSelected:
		detailsModel := hue_lights_details.InitModel(m.adapter, msg.Light)
		m.details = &detailsModel
		m.state = details
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
		}

	default:
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case list:
		return m.list.View()
	case details:
		return m.details.View()
	default:
		return ""
	}
}

func (m Model) findLight(bridgeId string, id string) *hue_control.Light {
	return m.adapter.Control.GetBridgeById(bridgeId).GetLightById(id)
}
