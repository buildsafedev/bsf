package search

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	bsfv1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/vulnerability"
)

type packageOptionModel struct {
	name     string
	version  string
	cursor   int
	errorMsg string
	selected map[string]bool
	vlm      versionListModel
	vulnResp *bsfv1.FetchVulnerabilitiesResponse
}

const (
	printableVulnCount = 10
)

var optionsShown bool

// InitOption initializes the option model
func initOption(name, version string, vlm versionListModel, vulnResp *bsfv1.FetchVulnerabilitiesResponse) *packageOptionModel {
	return &packageOptionModel{
		name:     name,
		version:  version,
		cursor:   0,
		selected: make(map[string]bool, 0),
		vlm:      vlm,
		vulnResp: vulnResp,
	}

}

// TODO: We'll need to support Buildtime dependencies in future.
var choices = []string{"Development", "Runtime"}

func (m packageOptionModel) Init() tea.Cmd {
	return nil
}

func (m packageOptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, KeyMap.Back):
			return m.vlm.Update(msg)

		case key.Matches(msg, KeyMap.Space):
			// Toggle selection
			if m.selected == nil {
				m.selected = make(map[string]bool, 0)
			}
			m.selected[choices[m.cursor]] = !m.selected[choices[m.cursor]]

		case key.Matches(msg, KeyMap.Up):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		case key.Matches(msg, KeyMap.Down):
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}
		case key.Matches(msg, KeyMap.Enter):
			if currentMode != modeOption {
				break
			}

			v := initVersionConstraints(choices[m.cursor], m.name, m.version)
			p := tea.NewProgram(v)
			if err := p.Start(); err != nil {
				m.errorMsg = fmt.Sprintf("Error starting version constraints model: %s", err.Error())
			}

			return m, tea.Quit

		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m packageOptionModel) View() string {
	if m.errorMsg != "" {
		// todo: debug how to render this properly. Maybe add a spiner Model/View?
		// also exist with status code 1 instead of 0
		return m.errorMsg
	}
	s := strings.Builder{}

	if m.vulnResp != nil {
		if len(m.vulnResp.Vulnerabilities) == 0 {
			s.WriteString(styles.SucessStyle.Render(fmt.Sprintf("%d vulnerabilities found for %s-%s", len(m.vulnResp.Vulnerabilities), m.name, m.version)))
		}
		if len(m.vulnResp.Vulnerabilities) != 0 {
			s.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("%d vulnerabilities found for %s-%s", len(m.vulnResp.Vulnerabilities), m.name, m.version)))
			s.WriteString("\n\n")
			s.WriteString(styles.TitleStyle.Render(fmt.Sprintf("%-20s %-10s %-10s %-20s ", "CVE", "Severity", "Score", "Vector")))
			s.WriteString("\n")

			sortedValues := vulnerability.SortVulnerabilities(m.vulnResp.Vulnerabilities)
			printableVulns := getTopVulnerabilities(sortedValues)

			for _, v := range printableVulns {
				// todo: maybe we should add style based on severity
				if v.Cvss == nil {
					continue
				}

				// let's pick the first cvss
				s.WriteString(styles.TextStyle.Render(fmt.Sprintf("%-20s %-10s %-10f %-20s", v.Id, v.Severity, v.Cvss[0].Metrics.BaseScore, vulnerability.DeriveAV(v.Cvss[0].Vector))))
				s.WriteString("\n")
			}

			if len(m.vulnResp.Vulnerabilities) > printableVulnCount {
				s.WriteString(styles.HintStyle.Render(fmt.Sprintf("More %d found, run `bsf scan %s:%s`", len(m.vulnResp.Vulnerabilities)-printableVulnCount, m.name, m.version)))
			}
		}
	}
	s.WriteString("\n\n")

	s.WriteString(styles.TitleStyle.Render(fmt.Sprintf("Where would you like %s-%s to be added?", m.name, m.version)))
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
	currentMode = modeOption
	return s.String()
}

func getTopVulnerabilities(allVuln []*bsfv1.Vulnerability) []*bsfv1.Vulnerability {
	if len(allVuln) <= printableVulnCount {
		return allVuln
	}
	return allVuln[:10]
}

func newConfFromSelectedPackages(name, version string, selected map[string]bool) hcl2nix.Packages {
	// since only package can searched and selected at a time, we can safely assume this.
	packages := hcl2nix.Packages{
		Development: make([]string, 0, 1),
		Runtime:     make([]string, 0, 1),
	}
	for _, c := range choices {
		if selected[c] {
			switch c {
			case "Development":
				packages.Development = append(packages.Development, name+"@"+version)
			case "Runtime":
				packages.Runtime = append(packages.Runtime, name+"@"+version)
			}
		}
	}
	return packages
}
