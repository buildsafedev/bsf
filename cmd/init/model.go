package init

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
	btemplate "github.com/buildsafedev/bsf/pkg/nix/template"
	"github.com/charmbracelet/bubbles/spinner"
)

var (
	// ANSI codes reference- https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
	textStyle    = styles.TextStyle.Render
	sucessStyle  = styles.SucessStyle.Render
	spinnerStyle = styles.SpinnerStyle
	helpStyle    = styles.HelpStyle.Render
	errorStyle   = styles.ErrorStyle.Render
	stages       = 4
)

type model struct {
	spinner  spinner.Model
	sc       *search.Client
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
		m.stageMsg = textStyle("Initializing project, detecting project language..  ")
		return nil
	case 1:
		// var err error
		pt, pd, err := langdetect.FindProjectType()
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}
		m.stageMsg = textStyle("Detected language as " + string(pt))
		if m.pt == langdetect.Unknown {
			m.permMsg = errorStyle("Project language isn't currently supported. Some features might not work.")
		}
		m.pt = pt
		m.pd = pd
		return nil
	case 2:
		m.stageMsg = textStyle("Resolving dependencies... ")
		return nil

	case 3:
		_, err := createBsfDirectory()
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		fh, err := NewFileHandlers()
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}
		defer fh.ModFile.Close()
		defer fh.LockFile.Close()
		defer fh.FlakeFile.Close()
		defer fh.DefFlakeFile.Close()

		conf := generatehcl2NixConf(m.pt, m.pd)
		err = hcl2nix.WriteConfig(conf, fh.ModFile)
		if err != nil {
			m.stageMsg = errorStyle("Is the project already initialised? \n bsf.hcl found. \n", err.Error())
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		defer cancel()

		allPackages, err := hcl2nix.ResolvePackages(ctx, m.sc, conf.Packages)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		err = hcl2nix.GenerateLockFile(allPackages, fh.LockFile)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		cr := hcl2nix.ResolveCategoryRevisions(conf.Packages, allPackages)
		err = btemplate.GenerateFlake(btemplate.Flake{
			// Description:         "bsf flake",
			NixPackageRevisions: cr.Revisions,
			DevPackages:         cr.Development,
			RuntimePackages:     cr.Runtime,
		}, fh.FlakeFile)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		// todo: there should be a generic method "GenerateApplicationModule" that can switch between different project types
		err = btemplate.GenerateGoModule(conf.GoModule, fh.DefFlakeFile)
		if err != nil {
			m.stageMsg = errorStyle(err.Error())
			return err
		}

		return nil
	}

	return nil
}

// FileHandlers holds file handlers
type FileHandlers struct {
	ModFile      *os.File
	LockFile     *os.File
	FlakeFile    *os.File
	DefFlakeFile *os.File
}

// NewFileHandlers creates new file handlers
func NewFileHandlers() (*FileHandlers, error) {
	modFile, err := GetOrCreateFile("bsf.hcl")
	if err != nil {
		return nil, err
	}

	lockFile, err := os.Create("bsf.lock")
	if err != nil {
		return nil, err
	}

	flakeFile, err := os.Create("bsf/flake.nix")
	if err != nil {
		return nil, err
	}

	defFlakeFile, err := os.Create("bsf/default.nix")
	if err != nil {
		return nil, err
	}

	return &FileHandlers{
		ModFile:      modFile,
		LockFile:     lockFile,
		FlakeFile:    flakeFile,
		DefFlakeFile: defFlakeFile,
	}, nil
}

// GetOrCreateFile gets or creates a file if it doesn't exist
func GetOrCreateFile(path string) (*os.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Create(path)
	}

	return os.Open(path)
}
