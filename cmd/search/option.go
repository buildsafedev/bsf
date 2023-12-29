package search

import (
	"fmt"
	"os"
	"strings"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/generate"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type packageOptionModel struct {
	item     list.Item
	cursor   int
	errorMsg string
	selected map[string]bool // Track selected options
	sc       *search.Client
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

			err = hcl2nix.AddPackages(data, newConfFromSelectedPackages(m.item, m.selected), fh.ModFile)
			if err != nil {
				m.errorMsg = fmt.Sprintf(errorStyle.Render("Error updating bsf.hcl: %s", err.Error()))
				return m, tea.Quit
			}

			err = generate.Generate(fh, m.sc)
			if err != nil {
				m.errorMsg = fmt.Sprintf(errorStyle.Render("Error generating files: %s", err.Error()))
				return m, tea.Quit
			}

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
	if m.errorMsg != "" {
		// todo: debug how to render this properly. Maybe add a spiner Model/View?
		// also exist with status code 1 instead of 0
		return m.errorMsg
	}
	s := strings.Builder{}

	pkg := m.item.(pkgVersionItem).pkg
	s.WriteString(styles.TitleStyle.Render(fmt.Sprintf("Where would you like %s-%s to be added?", pkg.Name, pkg.Version)))
	s.WriteString("\n\n")

	for i := 0; i < len(choices); i++ {
		if m.selected[choices[i]] {
			s.WriteString(styles.SelectedOptionStyle.Render("✔ " + choices[i]))
		} else if m.cursor == i {
			s.WriteString(styles.CursorOptionStyle.Render("-> " + choices[i]))
		} else {
			s.WriteString(styles.OptionStyle.Render("  " + choices[i]))
		}
		s.WriteString("\n")
	}

	s.WriteString(styles.HelpStyle.Render("\n(↑↓ to move cursor, space to select/unselect, enter to submit)\n"))

	return s.String()
}

func newConfFromSelectedPackages(item list.Item, selected map[string]bool) hcl2nix.Packages {
	pkg := item.(pkgVersionItem).pkg

	// since only package can searched and selected at a time, we can safely assume this.
	packages := hcl2nix.Packages{
		Development: make([]string, 0, 1),
		Runtime:     make([]string, 0, 1),
	}
	for _, c := range choices {
		if selected[c] {
			switch c {
			case "Development":
				packages.Development = append(packages.Development, pkg.Name+"@"+pkg.Version)
			case "Runtime":
				packages.Runtime = append(packages.Runtime, pkg.Name+"@"+pkg.Version)
			}
		}
	}
	return packages
}
