package search

import (
	"fmt"
	"os"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/generate"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type versionConstraintsModel struct {
	choices             []string
	cursor              int
	selectedConstraints string
	selected            map[string]bool
	context             string
	errorMsg            string
	name                string
	version             string
}

func (m *versionConstraintsModel) Init() tea.Cmd {
	return nil
}

func initVersionConstraints(context, name, version string) *versionConstraintsModel {
	choices := []string{"pinned version", "allow minor version updates", "allow patch version updates"}
	selected := make(map[string]bool)
	return &versionConstraintsModel{
		choices:  choices,
		cursor:   0,
		selected: selected,
		context:  context,
		name:     name,
		version:  version,
	}
}

func (m *versionConstraintsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, KeyMap.Down):
			m.cursor++
			if m.cursor > len(m.choices)-1 {
				m.cursor = 0
			}
		case key.Matches(msg, KeyMap.Up):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		case key.Matches(msg, KeyMap.Enter):

			fh, err := hcl2nix.NewFileHandlers(true)
			if err != nil {
				m.errorMsg = fmt.Sprintf("Error creating file handlers: %s", err.Error())
				return m, tea.Quit
			}

			data, err := os.ReadFile("bsf.hcl")
			if err != nil {
				m.errorMsg = fmt.Sprintf(errorStyle.Render("Error reading bsf.hcl: %s", err.Error()))
				return m, tea.Quit
			}

			// changing file handler to allow writes
			fh.ModFile, err = os.Create("bsf.hcl")
			if err != nil {
				m.errorMsg = fmt.Sprintf(errorStyle.Render("Error creating bsf.hcl: %s", err.Error()))
				return m, tea.Quit
			}

			err = hcl2nix.AddPackages(data, newConfFromSelectedPackages(m.name, m.version, m.selected), fh.ModFile)
			if err != nil {
				m.errorMsg = fmt.Sprintf(errorStyle.Render("Error updating bsf.hcl: %s", err.Error()))
				return m, tea.Quit
			}

			err = generate.Generate(fh, sc)
			if err != nil {
				m.errorMsg = fmt.Sprintf(errorStyle.Render("Error generating files: %s", err.Error()))
				return m, tea.Quit
			}

			return m, tea.Quit
		case key.Matches(msg, KeyMap.Space):
			choice := m.choices[m.cursor]
			m.selected[choice] = !m.selected[choice]
		}
	}
	return m, nil
}

func (m versionConstraintsModel) View() string {
	var s strings.Builder

	for i, choice := range m.choices {
		if m.selected[choice] {
			s.WriteString(styles.SelectedOptionStyle.Render("  [x] " + choice))
		} else if m.cursor == i {
			s.WriteString(styles.BaseStyle.Render(" ->  " + choice))
		} else {
			s.WriteString(styles.BaseStyle.Render(" []" + choice))
		}
		s.WriteString("\n")
	}
	return s.String()
}
