package search

import (
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type pkgitem struct {
	pkg search.Package
}

func (i pkgitem) Title() string       { return i.pkg.Name }
func (i pkgitem) Description() string { return i.pkg.Version }
func (i pkgitem) FilterValue() string { return i.pkg.Version }

type pkgListModel struct {
	pkgList            list.Model
	packageOptionModel packageOptionModel
}

func (m pkgListModel) Init() tea.Cmd {
	return nil
}

func (m pkgListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			item := m.pkgList.SelectedItem()
			m.packageOptionModel.item = item
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.pkgList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.pkgList, cmd = m.pkgList.Update(msg)
	return m, cmd
}

func (m pkgListModel) View() string {
	if m.packageOptionModel.item != nil {
		return m.packageOptionModel.View()
	}
	return docStyle.Render(m.pkgList.View())
}
