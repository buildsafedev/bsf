package template

import (
	"html/template"
	"io"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

const (
	rustTmpl = `
	{pkgs,rustPkgs}:
 	  (rustPkgs pkgs).workspace.{{ .CrateName }} {}
    `
)

type RustApp struct {
	CrateName           string
	RustVersion         string
	RustToolChain       string
	RustChannel         string
	RustProfile         string
	ExtraRustComponents []string
	Release             bool
}

// GenerateRustApp generates default flake
func GenerateRustApp(fl *hcl2nix.RustApp, wr io.Writer) error {
	data := RustApp{
		CrateName:           fl.CrateName,
		RustVersion:         fl.RustVersion,
		RustToolChain:       fl.RustToolChain,
		RustChannel:         fl.RustChannel,
		RustProfile:         fl.RustProfile,
		ExtraRustComponents: fl.ExtraRustComponents,
		Release:             fl.Release,
	}
	t, err := template.New("rust").Parse(rustTmpl)
	if err != nil {
		return err
	}

	err = t.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}
