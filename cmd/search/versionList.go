package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
)

const (
	totalCols = 6
)

type versionListModel struct {
	pkgList      list.Model
	versionTable table.Model
}

func convFPR2Rows(versions *buildsafev1.FetchPackagesResponse) []table.Row {
	items := make([]table.Row, 0, len(versions.Packages))

	sortedPackages := search.SortPackages(versions.Packages)
	for _, pkg := range sortedPackages {
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
			time.Unix(int64(pkg.EpochSeconds), 0).Format("2006-01-02 15:04"),
		})
	}
	return items
}

// initVersionTable initializes the version table
func initVersionTable(searchList list.Model, versions *buildsafev1.FetchPackagesResponse) *versionListModel {
	columns := []table.Column{
		{Title: "Name", Width: frameWidth / totalCols},
		{Title: "Version", Width: frameWidth / totalCols},
		{Title: "License", Width: frameWidth / totalCols},
		{Title: "Free", Width: frameWidth / totalCols},
		{Title: "Homepage", Width: frameWidth * 2 / totalCols},
		{Title: "Date", Width: frameWidth / totalCols},
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
			if len(row) != totalCols {
				// TODO: return errMsg
				return m, tea.Quit
			}
			name := row[0]
			version := row[1]
			vulnResp, err := sc.FetchVulnerabilities(context.Background(), &buildsafev1.FetchVulnerabilitiesRequest{
				Name:    name,
				Version: version,
			})
			if err != nil {
				fmt.Println(errorStyle.Render(err.Error()))
			}

			return initOption(name, version, m, vulnResp).Update(msg)
		}
	}

	m.versionTable, cmd = m.versionTable.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m versionListModel) View() string {
	s := strings.Builder{}

	s.WriteString(styles.BaseStyle.Render(m.versionTable.View() + "\n"))
	s.WriteString(styles.HelpStyle.Render("\n(↑↓ to move cursor,  enter to submit, esc to previous prompt, ctr+c to quit)\n"))
	currentMode = modeVersion
	return s.String()
}
