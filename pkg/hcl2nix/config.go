package hcl2nix

import (
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"

	bstrings "github.com/buildsafedev/bsf/pkg/strings"
)

// Config for hcl2nix
type Config struct {
	Packages  Packages       `hcl:"packages,block"`
	GoModule  *GoModule      `hcl:"gomodule,block"`
	PoetryApp *PoetryApp     `hcl:"poetryapp,block"`
	RustApp   *RustApp       `hcl:"rustapp,block"`
	Export    []ExportConfig `hcl:"export,block"`
}

// Packages holds package parameters
type Packages struct {
	// Maybe these should be of type Set? https://github.com/deckarep/golang-set
	Development []string `hcl:"development"`
	Runtime     []string `hcl:"runtime"`
}

// WriteConfig writes packages to writer
func WriteConfig(config Config, wr io.Writer) error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&config, f.Body())
	_, err := f.WriteTo(wr)
	if err != nil {
		return err
	}
	return nil
}

// ReadConfig reads config from bytes and returns Config. If any errors are encountered, they are written to dstErr
func ReadConfig(src []byte, dstErr io.Writer) (*Config, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL(src, "bsf.hcl")
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			dstErr,
			parser.Files(),
			78,
			true,
		)
		wr.WriteDiagnostics(diags)
		return nil, diags
	}

	var config Config
	diags = gohcl.DecodeBody(f.Body, nil, &config)
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			dstErr,
			parser.Files(),
			78,
			true,
		)
		wr.WriteDiagnostics(diags)
		return nil, diags
	}

	return &config, nil
}

// AddPackages updates config with new packages. It appends new packages to existing packages
func AddPackages(src []byte, packages Packages, wr io.Writer) error {
	existingConfig, err := ReadConfig(src, io.Discard)
	if err != nil {
		return err
	}

	// append new packages to existing packages
	existingConfig.Packages.Development = bstrings.SliceToSet(append(existingConfig.Packages.Development, packages.Development...))
	existingConfig.Packages.Runtime = bstrings.SliceToSet(append(existingConfig.Packages.Runtime, packages.Runtime...))

	err = WriteConfig(*existingConfig, wr)
	if err != nil {
		return err
	}

	return nil
}

// SetPackages updates config with new packages. It replaces existing packages with new packages
func SetPackages(src []byte, packages Packages, wr io.Writer) error {
	existingConfig, err := ReadConfig(src, io.Discard)
	if err != nil {
		return err
	}

	existingConfig.Packages = packages

	err = WriteConfig(*existingConfig, wr)
	if err != nil {
		return err
	}

	return nil

}
