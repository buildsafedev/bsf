package search

import (
	"fmt"
	"os"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/mod/semver"
)

type versionConstraintsModel struct {
	choices             []string
	cursor              int
	constraint          map[string]bool
	errorMsg            string
	name                string
	version             string
	selectedConstraints string
	env                 map[string]bool
	pkgoption           packageOptionModel
}

func (m *versionConstraintsModel) Init() tea.Cmd {
	return nil
}

func initVersionConstraints(name, version string, env map[string]bool, pkgoption packageOptionModel) *versionConstraintsModel {
	var choices []string
	var constriant map[string]bool
	if semver.IsValid("v" + version) {
		choices = []string{"pinned version", "allow minor version updates", "allow patch version updates"}
		constriant = map[string]bool{"allow patch version updates": true}
	} else {
		constriant = map[string]bool{"pinned version": true}
	}

	if !env["Development"] && !env["Runtime"] {
		env["Development"] = true
	}

	return &versionConstraintsModel{
		choices:    choices,
		cursor:     2,
		constraint: constriant,
		name:       name,
		version:    version,
		env:        env,
		pkgoption:  pkgoption,
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
		case key.Matches(msg, KeyMap.Back):
			return m.pkgoption.Update(msg)
		case key.Matches(msg, KeyMap.Space):
			if m.constraint[m.choices[m.cursor]] {
				m.constraint[m.choices[m.cursor]] = false
			} else {
				for k := range m.constraint {
					m.constraint[k] = false
				}
				m.constraint[m.choices[m.cursor]] = true
			}

		case key.Matches(msg, KeyMap.Enter):
			return m.updateVersionConstraint()
		}
	}
	return m, nil
}

func (m versionConstraintsModel) View() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString(styles.TitleStyle.Render(fmt.Sprintf("What type of updates would you like for %s-%s?", m.name, m.version)))
	s.WriteString("\n\n")
	for i, choice := range m.choices {
		if m.constraint[choice] {
			s.WriteString(styles.SelectedOptionStyle.Render("✔ " + choice))
		} else if m.cursor == i {
			s.WriteString(styles.CursorOptionStyle.Render("-> " + choice))
		} else {
			s.WriteString(styles.OptionStyle.Render("  " + choice))
		}
		s.WriteString("\n")

	}
	s.WriteString(styles.HelpStyle.Render("\n(↑↓ to move cursor, space to select/unselect, enter to submit)\n"))

	return s.String()
}

func (m *versionConstraintsModel) updateVersionConstraint() (tea.Model, tea.Cmd) {
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

	switch {
	case m.constraint["allow minor version updates"]:
		m.selectedConstraints = "^"
	case m.constraint["allow patch version updates"]:
		m.selectedConstraints = "~"
	}
	if m.selectedConstraints == "" && !m.constraint["pinned version"] {
		m.selectedConstraints = "~"
	}

	// changing file handler to allow writes
	fh.ModFile, err = os.Create("bsf.hcl")
	if err != nil {
		m.errorMsg = fmt.Sprintf(errorStyle.Render("Error creating bsf.hcl: %s", err.Error()))
		return m, tea.Quit
	}

	err = hcl2nix.AddPackages(data, newConfFromSelectedPackages(m.name, m.version, m.selectedConstraints, m.env), fh.ModFile)
	if err != nil {
		m.errorMsg = fmt.Sprintf(errorStyle.Render("Error updating bsf.hcl: %s", err.Error()))
		return m, tea.Quit
	}

	return m, tea.Quit
}
