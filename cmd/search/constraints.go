package search

import (
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type versionConstraintsModel struct {
	choices             []string
	cursor              int
	selectedConstraints string
	selected            map[string]bool
	context             string
}

func (m *versionConstraintsModel) Init() tea.Cmd {
	return nil
}

func initVersionConstraints(context string) *versionConstraintsModel {
	choices := []string{"pinned version", "allow minor version updates", "allow patch version updates"}
	selected := make(map[string]bool)
	return &versionConstraintsModel{
		choices:  choices,
		cursor:   0,
		selected: selected,
		context:  context,
	}
}

func (m *versionConstraintsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "j", "down":
			m.cursor++
			if m.cursor > len(m.choices)-1 {
				m.cursor = 0
			}
		case "k", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		case "enter":
			m.selectedConstraints = m.choices[m.cursor]
		}
	}
	return m, nil
}

func (m versionConstraintsModel) View() string {
	var s strings.Builder

	for i, choice := range m.choices {
		if m.selected[m.choices[i]] {
			s.WriteString(styles.SelectedOptionStyle.Render("✔ " + choice))
		} else if m.cursor == i {
			s.WriteString(styles.CursorOptionStyle.Render("-> " + choice))
		} else {
			s.WriteString(styles.OptionStyle.Render("  " + choice))
		}
		s.WriteString("\n")
	}
	return s.String()
}
