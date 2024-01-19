package dc_textinput

import (
	"fmt"

	"github.com/remshams/device-control/tui/styles"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Edit    key.Binding
	Discard key.Binding
	Apply   key.Binding
}

var TextInputKeyMap = KeyMap{
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Discard: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "discard"),
	),
	Apply: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "apply"),
	),
}

type Model struct {
	label string
	unit  string
	Input textinput.Model
}

func New(label string, unit string) Model {
	return Model{
		label: label,
		unit:  unit,
		Input: CreateTextInputModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, TextInputKeyMap.Edit):
			m.Input.Focus()
		case key.Matches(msg, TextInputKeyMap.Discard):
			m.Input.Blur()
		case key.Matches(msg, TextInputKeyMap.Apply):
			m.Input.Blur()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return CreateTextInputView(m.Input, m.label, m.unit)
}

func CreateTextInputModel() textinput.Model {
	model := textinput.New()
	model.TextStyle = styles.TextAccentColor
	return model
}

func CreateTextInputView(model textinput.Model, label string, unit string) string {
	return fmt.Sprintf("%s %s%s", label, model.View(), styles.TextAccentColor.Render(unit))
}
