package hcl2nix

import (
	"io"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"

	bstrings "github.com/buildsafedev/bsf/pkg/strings"
)

// Config for hcl2nix
type Config struct {
	Packages Packages  `hcl:"packages,block"`
	GoModule *GoModule `hcl:"gomodule,block"`
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

// ReadConfig reads config from bytes and returns Config
func ReadConfig(src []byte) (*Config, error) {
	parser := hclparse.NewParser()
	f, err := parser.ParseHCL(src, "bsf.mod")
	if err != nil {
		return nil, err
	}

	var config Config
	diags := gohcl.DecodeBody(f.Body, nil, &config)
	if diags.HasErrors() {
		return nil, diags
	}
	return &config, nil
}

// AddPackages updates config with new packages. It appends new packages to existing packages
func AddPackages(config Config, src []byte, wr io.Writer) error {
	existingConfig, err := ReadConfig(src)
	if err != nil {
		return err
	}

	// append new packages to existing packages
	existingConfig.Packages.Development = bstrings.SliceToSet(append(existingConfig.Packages.Development, config.Packages.Development...))
	existingConfig.Packages.Runtime = bstrings.SliceToSet(append(existingConfig.Packages.Runtime, config.Packages.Runtime...))

	err = WriteConfig(*existingConfig, wr)
	if err != nil {
		return err
	}

	return nil
}
