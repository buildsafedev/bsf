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
	WorkspaceSrc                  string
	CrateName                     string
	RustVersion                   string
	RustToolChain                 string
	RustChannel                   string
	RustProfile                   string
	ExtraRustComponents           []string
	Release                       bool
	RootFeatures                  []string
	FetchCrateAlternativeRegistry string
	HostPlatformCpu               string
	HostPlatformFeatures          []string
	CargoUnstableFlags            []string
	RustcLinkFlags                []string
	RustcBuildFlags               []string
}

// GenerateRustApp generates default flake
func GenerateRustApp(fl *hcl2nix.RustApp, wr io.Writer) error {
	data := RustApp{
		WorkspaceSrc:                  fl.WorkspaceSrc,
		CrateName:                     fl.CrateName,
		RustVersion:                   fl.RustVersion,
		RustToolChain:                 fl.RustToolChain,
		RustChannel:                   fl.RustChannel,
		RustProfile:                   fl.RustProfile,
		ExtraRustComponents:           fl.ExtraRustComponents,
		Release:                       fl.Release,
		RootFeatures:                  fl.RootFeatures,
		FetchCrateAlternativeRegistry: fl.FetchCrateAlternativeRegistry,
		HostPlatformCpu:               fl.HostPlatformCpu,
		HostPlatformFeatures:          fl.HostPlatformFeatures,
		CargoUnstableFlags:            fl.CargoUnstableFlags,
		RustcLinkFlags:                fl.RustcLinkFlags,
		RustcBuildFlags:               fl.RustcBuildFlags,
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
