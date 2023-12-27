package search

import (
	"fmt"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type packageOptionModel struct {
	item     list.Item
	cursor   int
	selected map[string]bool // Track selected options
}

// Maybe Buildtime should be an option?
var choices = []string{"Development", "Runtime"}

func (m packageOptionModel) Init() tea.Cmd {
	return nil
}

func (m packageOptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case " ": // spacebar
			// Toggle selection
			if m.selected == nil {
				m.selected = make(map[string]bool, 0)
			}
			m.selected[choices[m.cursor]] = !m.selected[choices[m.cursor]]

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "enter":
			// TODO: Handle Enter key
			return m, tea.Quit

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m packageOptionModel) View() string {
	s := strings.Builder{}

	pkg := m.item.(pkgVersionItem).pkg
	s.WriteString(styles.TitleStyle.Render(fmt.Sprintf("Where would you like %s-%s to be added?", pkg.Name, pkg.Version)))
	s.WriteString("\n\n")

	for i := 0; i < len(choices); i++ {
		if m.selected[choices[i]] {
			s.WriteString(styles.SelectedOptionStyle.Render("(•) " + choices[i]))
		} else if m.cursor == i {
			s.WriteString(styles.SelectedOptionStyle.Render("(•) " + choices[i]))
		} else {
			s.WriteString(styles.OptionStyle.Render("( ) " + choices[i]))
		}
		s.WriteString("\n")
	}

	s.WriteString(styles.HelpStyle.Render("\n(•press q to quit  •press space to select/unselect •press enter to submit)\n"))

	return s.String()
}
