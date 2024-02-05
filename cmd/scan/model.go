package scan

import (
	"fmt"
	"strings"

	"github.com/buildsafedev/bsf/cmd/search"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	bsfv1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
)

var optionsShown bool

type vulnListModel struct {
	vulnTable table.Model
}

func convVulns2Rows(vulnerabilities *bsfv1.FetchVulnerabilitiesResponse) []table.Row {
	items := make([]table.Row, 0, len(vulnerabilities.Vulnerabilities))

	sortedVulns := search.SortVulnerabilities(vulnerabilities.Vulnerabilities)
	for _, v := range sortedVulns {
		items = append(items, table.Row{
			v.Id,
			v.Severity,
			fmt.Sprint(v.Cvss[0].Metrics.BaseScore),
			search.DeriveAV(v.Cvss[0].Vector),
		})
	}
	return items
}

func initVulnTable(vulnResp *bsfv1.FetchVulnerabilitiesResponse) *vulnListModel {

	frameHeight, frameWidth := styles.DocStyle.GetFrameSize()

	cols := 6
	columns := []table.Column{
		{Title: "CVE", Width: frameWidth / cols},
		{Title: "Severity", Width: frameWidth / cols},
		{Title: "Score", Width: frameWidth / cols},
		{Title: "Vector", Width: frameWidth / cols},
	}

	rows := convVulns2Rows(vulnResp)
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

	return &vulnListModel{
		vulnTable: t,
	}

}

func (m vulnListModel) Init() tea.Cmd {
	return nil
}

// Update handles events and updates the model accordingly
func (m vulnListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, search.KeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, search.KeyMap.Back):
			// without this, we would go back to the search model
			if optionsShown {
				break
			}
			return initVulnTable(m.vulnTable.Items()), nil
		}
	}
}

// View renders the user interface based on the current model
func (m vulnListModel) View() string {
	s := strings.Builder{}

	// Header
	s.WriteString(styles.BaseStyle.Render(m.vulnTable.View() + "\n"))
	s.WriteString(styles.HelpStyle.Render("\n(↑↓ to move cursor,  esc to previous prompt, ctr+c to quit)\n"))
	return s.String()
}
