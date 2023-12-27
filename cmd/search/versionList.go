package search

import (
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type pkgVersionItem struct {
	pkg search.Package
}

func (i pkgVersionItem) Title() string       { return i.pkg.Name }
func (i pkgVersionItem) Description() string { return i.pkg.Version }
func (i pkgVersionItem) FilterValue() string { return i.pkg.Version }

type versionListModel struct {
	// todo: maybe this should be a table?
	versionList        list.Model
	packageOptionModel packageOptionModel
}

func (m versionListModel) Init() tea.Cmd {
	return nil
}

func (m versionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.packageOptionModel.item != nil {
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
			m.packageOptionModel.item = item
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.versionList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.versionList, cmd = m.versionList.Update(msg)
	return m, cmd
}

func (m versionListModel) View() string {
	if m.packageOptionModel.item != nil {
		return m.packageOptionModel.View()
	}
	return docStyle.Render(m.versionList.View())
}
