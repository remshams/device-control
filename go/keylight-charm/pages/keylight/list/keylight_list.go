package keylight_list

import (
	"keylight-charm/keylight"
	"keylight-control/control"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewState string

type Model struct {
	keylights       []control.Keylight
	keylightAdapter *keylight.KeylightAdapter
	table           table.Model
}

func InitModel(keylightAdapter *keylight.KeylightAdapter, keylights []control.Keylight) Model {
	model := Model{keylights: []control.Keylight{}, keylightAdapter: keylightAdapter, table: createTable(keylights)}
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmd = tea.Quit
		case "enter":
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

func createTable(keylights []control.Keylight) table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Name", Width: 30},
	}
	rows := []table.Row{}
	for _, keylight := range keylights {
		rows = append(rows, table.Row{strconv.Itoa(keylight.Metadata.Id), keylight.Metadata.Name})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
