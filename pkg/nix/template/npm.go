package template

import (
	"html/template"
	"io"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

const (
	npmTmpl = `
    { stdenv, buildNodeModules, lib, nodejs, npmHooks }:
	  stdenv.mkDerivation {
	  pname = "{{ .PackageName }}";
	  version = "0.1.0";
	  src = ../.;

	  nativeBuildInputs = [
		  buildNodeModules.hooks.npmConfigHook
		  nodejs
		  npmHooks.npmInstallHook
	  ];

	  nodeModules = buildNodeModules.fetchNodeModules {
		  packageRoot = {{ .PackageRoot }};
	  };
	}
    `
)

type npmApp struct {
	// PackageName: Name of the package.
	PackageName string
	// PackageJsonPath: Source path to the package.json and package-lock.json file.
	PackageRoot string
}

// GenerateNpmApp generates default flake
func GenerateNpmApp(fl *hcl2nix.JsNpmApp, wr io.Writer) error {
	data := npmApp{
		PackageName: fl.PackageName,
		PackageRoot: parentFolder(fl.PackageRoot),
	}

	t, err := template.New("npm").Parse(npmTmpl)
	if err != nil {
		return err
	}

	err = t.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}
