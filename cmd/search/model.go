package search

import (
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	pkg search.Package
}

func (i item) Title() string       { return i.pkg.Name }
func (i item) Description() string { return i.pkg.Version }
func (i item) FilterValue() string { return i.pkg.Version }

type model struct {
	// todo: maybe this should be a table?
	versionList        list.Model
	selected           bool
	packageOptionModel packageOptionModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.selected {
		nm, cmd := m.packageOptionModel.Update(msg)
		return nm, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			item := m.versionList.SelectedItem()
			_ = item
			m.selected = true
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.versionList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.versionList, cmd = m.versionList.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.selected {
		return m.packageOptionModel.View()
	}
	return docStyle.Render(m.versionList.View())
}
