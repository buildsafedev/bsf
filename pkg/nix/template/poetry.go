package template

import (
	"html/template"
	"io"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

const (
	poetryTmpl = `
    { pkgs, mkPoetryApplication }:
     let app = mkPoetryApplication {
            projectDir = {{.ProjectDir}};
            src = {{.Src}};
            pyproject = {{.Pyproject}};
            poetryLock = {{.Poetrylock}};
            python = pkgs.python3;
            preferWheels = {{.PreferWheels}};
			{{ if gt (len .CheckGroups) 0}}
			checkGroups = [
		      {{ range $value := .CheckGroups }}"{{ $value }}" {{ end }}
	       ];
		   {{ end }}
    };
	in app.dependencyEnv
    `
)

type poetryApp struct {
	// ProjectDir: path to the root of the project.
	ProjectDir string
	// Src: project source (defaults to projectDir).
	Src string
	// Pyproject: path to pyproject.toml (defaults to projectDir + "/pyproject.toml").
	Pyproject string
	// Poetrylock: poetry.lock file path (defaults to projectDir + "/poetry.lock").
	Poetrylock string
	// PreferWheels: Use wheels rather than sdist as much as possible (defaults to false).
	PreferWheels bool
	// CheckGroups: Which Poetry 1.2.0+ dependency groups to run unit tests (defaults to [ "dev" ]).
	CheckGroups []string
}

// GeneratePoetryApp generates default flake
func GeneratePoetryApp(fl *hcl2nix.PoetryApp, wr io.Writer) error {
	data := poetryApp{
		ProjectDir:   parentFolder(fl.ProjectDir),
		Src:          parentFolder(fl.Src),
		Pyproject:    parentFolder(fl.Pyproject),
		Poetrylock:   parentFolder(fl.Poetrylock),
		PreferWheels: fl.PreferWheels,
		CheckGroups:  fl.CheckGroups,
	}

	t, err := template.New("python-poetry").Parse(poetryTmpl)
	if err != nil {
		return err
	}

	err = t.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}
