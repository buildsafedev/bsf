package search

import (
	"strings"

	buildsafev1 "github.com/buildsafedev/cloud-api/apis/v1"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
)

type versionListModel struct {
	pkgList      list.Model
	versionTable table.Model
}

func convFPR2Rows(versions *buildsafev1.FetchPackagesResponse) []table.Row {
	items := make([]table.Row, 0, len(versions.Packages))

	soredPackages := search.SortPackagesWithTimestamp(versions.Packages)
	for _, pkg := range soredPackages {
		free := "false"
		if pkg.Free {
			free = "true"
		}
		items = append(items, table.Row{
			pkg.Name,
			pkg.Version,
			pkg.SpdxId,
			free,
			pkg.Homepage,
		})
	}
	return items
}

// initVersionTable initializes the version table
func initVersionTable(pkgName string, searchList list.Model, versions *buildsafev1.FetchPackagesResponse) *versionListModel {
	cols := 6
	columns := []table.Column{
		{Title: "Name", Width: frameWidth / cols},
		{Title: "Version", Width: frameWidth / cols},
		{Title: "License", Width: frameWidth / cols},
		{Title: "Free", Width: frameWidth / cols},
		{Title: "Homepage", Width: frameWidth * 2 / cols},
	}

	rows := convFPR2Rows(versions)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(frameHeight*8/10),
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

	return &versionListModel{
		versionTable: t,
		pkgList:      searchList,
	}

}

func (m versionListModel) Init() tea.Cmd {
	return nil
}

func (m versionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.DocStyle.GetFrameSize()
		m.versionTable.SetWidth(msg.Width - h)
		m.versionTable.SetHeight(msg.Height - v)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, KeyMap.Back):
			// without this, we would go back to the search model
			if optionsShown {
				break
			}
			return InitSearch(m.pkgList.Items()), nil

		case key.Matches(msg, KeyMap.Enter):
			if currentMode != modeVersion {
				break
			}
			row := m.versionTable.SelectedRow()
			if len(row) != 5 {
				// TODO: return errMsg
				return m, tea.Quit
			}
			name := row[0]
			version := row[1]
			return initOption(name, version, m).Update(msg)
		}
	}

	m.versionTable, cmd = m.versionTable.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m versionListModel) View() string {
	s := strings.Builder{}

	s.WriteString(styles.BaseStyle.Render(m.versionTable.View() + "\n"))
	s.WriteString(styles.HelpStyle.Render("\n(↑↓ to move cursor,  enter to submit, esc to previous prompt,ctr+c to quit)\n"))
	currentMode = modeVersion
	return s.String()
}
