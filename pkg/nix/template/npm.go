package template

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

const (
	npmTmpl = `
    { stdenv, buildNodeModules, lib, nodejs, npmHooks }:
	  let
	    packageRoot = {{ .PackageRoot }};
		{{ if ne .PackageJSONPath ""}}
	    package = lib.importJSON {{ .PackageJSONPath }}; {{ end }}
		{{ if ne .PackageLockPath ""}}
		packageLock = lib.importJSON {{ .PackageLockPath }}; {{ end }}
      in
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
		inherit packageRoot {{ if ne .PackageJSONPath ""}}package{{ end }} {{ if ne .PackageLockPath ""}}packageLock{{ end }} ;
	  };
	}
    `
)

type npmApp struct {
	// PackageName: Name of the package.
	PackageName string
	// PackageJsonPath: Source path to the package.json and package-lock.json file.
	PackageRoot string
	// PackageJSONPath: Path to package.json file.
	PackageJSONPath string
	// PackageLockPath: Path to package-lock.json file.
	PackageLockPath string
}

// GenerateNpmApp generates default flake
func GenerateNpmApp(fl *hcl2nix.JsNpmApp, wr io.Writer) error {
	data := npmApp{
		PackageName: fl.PackageName,
		PackageRoot: parentFolder(fl.PackageRoot),
	}

	if fl.PackageJSONPath != "" {
		data.PackageJSONPath = modifyPath(parentFolder(fl.PackageJSONPath), "packagejson")
	}

	if fl.PackageLockPath != "" {
		data.PackageLockPath = modifyPath(parentFolder(fl.PackageLockPath), "packagelock")
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

func modifyPath(directory, fileType string) string {
	directory = strings.TrimSuffix(directory, ".")

	var newFileName string
	switch fileType {
	case "packagejson":
		newFileName = "package.json"
	case "packagelock":
		newFileName = "package-lock.json"
	}
	newPath := filepath.Join(directory, newFileName)
	return newPath
}
