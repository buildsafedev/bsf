package init

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/generate"
	bgit "github.com/buildsafedev/bsf/pkg/git"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
	"github.com/buildsafedev/bsf/pkg/nix/cmd"
)

var (
	// ANSI codes reference- https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
	textStyle    = styles.TextStyle.Render
	sucessStyle  = styles.SucessStyle.Render
	spinnerStyle = styles.SpinnerStyle
	// helpStyle    = styles.HelpStyle.Render
	errorStyle = styles.ErrorStyle.Render
	stages     = 4
)

type model struct {
	spinner  spinner.Model
	sc       buildsafev1.SearchServiceClient
	stageMsg string
	permMsg  string
	stage    int
	pt       langdetect.ProjectType
	pd       *langdetect.ProjectDetails
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
			cleanUp()
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
		cleanUp()
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
		m.stageMsg = textStyle("Initializing project, detecting project language..  ")
		return nil
	case 1:
		// var err error
		pt, _, err := langdetect.FindProjectType()
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}
		m.stageMsg = textStyle("Detected language as " + string(pt))
		if m.pt == langdetect.Unknown {
			m.permMsg = errorStyle("Project language isn't currently supported. Some features might not work.")
		}
		m.pt = pt
		return nil
	case 2:
		m.stageMsg = textStyle("Resolving dependencies... ")
		return nil

	case 3:
		fh, err := hcl2nix.NewFileHandlers(false)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}
		defer fh.ModFile.Close()
		defer fh.LockFile.Close()
		defer fh.FlakeFile.Close()
		defer fh.DefFlakeFile.Close()
		conf, err := generatehcl2NixConf(m.pt, m.pd)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}
		err = hcl2nix.WriteConfig(conf, fh.ModFile)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		err = generate.Generate(fh, m.sc)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		err = cmd.Lock()
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		err = bgit.Add("bsf/")
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		err = bgit.Ignore("bsf-result/")
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		return nil
	}

	return nil
}
