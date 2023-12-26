package init

import (
	"fmt"
	"time"

	"github.com/buildsafedev/bsf/pkg/langdetect"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// ANSI codes reference- https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("254")).Render
	sucessStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render
	stages       = 4
)

type model struct {
	spinner  spinner.Model
	stageMsg string
	permMsg  string
	stage    int
	pt       langdetect.ProjectType
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return spinner.TickMsg(spinner.TickMsg{Time: t})
	})
}

func (m *model) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinner.Points
}

func (m model) View() (s string) {
	s += fmt.Sprintf("\n %s%s%s\n\n", m.spinner.View(), " ", m.stageMsg)
	if m.permMsg != "" {
		s += fmt.Sprintf("\n %s\n", m.permMsg)
	}
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	var err error
	if m.stage >= stages {
		m.stageMsg = sucessStyle("Initialised sucessfully!")
		return m, tea.Quit
	}
	err = m.processStages(m.stage)
	if err != nil {
		return m, tea.Quit
	}
	m.stage++

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(spinner.TickMsg(spinner.TickMsg{Time: time.Now()}))
	return m, cmd
}

func (m *model) processStages(stage int) error {
	switch stage {
	case 0:
		m.stageMsg = textStyle("Initializing project  ")
		return nil
	case 1:
		_, err := createBsfDirectory()
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}
		return nil

	case 2:
		// var err error
		pt, err := langdetect.FindProjectType()
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}
		m.stageMsg = textStyle("Detected language as " + string(pt))
		m.pt = pt
		return nil
	case 3:
		if m.pt == langdetect.Unknown {
			m.permMsg = errorStyle("Project language/package manager isn't currently supported. Some features might not work.")
		}
		// _, err := generateNixFiles()
		m.stageMsg = textStyle("Generating flake.." + string(m.pt))
		return nil
	}

	return nil
}
