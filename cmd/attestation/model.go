package attestation

import (
	"fmt"
	"os"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type psModel struct {
	psTable table.Model
}

func convPredSubtoRows(psMap map[string][]string) []table.Row {
	items := make([]table.Row, 0, len(psMap))

	for pred, subjects := range psMap {
		for _, sub := range subjects {
			items = append(items, table.Row{
				pred,
				sub,
			})
		}

	}
	return items
}

func initPredSubTable(psMap map[string][]string) *psModel {

	frameWidth, frameHeight, err := term.GetSize(0)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	cols := 2
	columns := []table.Column{
		{Title: "Predicate", Width: frameWidth / 5 * cols},
		{Title: "Subject", Width: frameWidth / 5 * cols},
	}

	rows := convPredSubtoRows(psMap)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(frameHeight*6/10),
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

	return &psModel{
		psTable: t,
	}

}

func (m psModel) Init() tea.Cmd {
	return nil
}

// Update handles events and updates the model accordingly
func (m psModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.psTable, cmd = m.psTable.Update(msg)
	return m, cmd
}

// View renders the user interface based on the current model
func (m psModel) View() string {
	s := strings.Builder{}

	// Header
	s.WriteString(styles.BaseStyle.Render(m.psTable.View() + "\n"))
	s.WriteString(styles.HelpStyle.Render("\n(↑↓ to move cursor, ctr+c to quit)\n"))
	return s.String()
}
