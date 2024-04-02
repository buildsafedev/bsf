package template

import (
	"html/template"
	"io"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

const (
	rustTmpl = `
	{ pkgs}:
    rustPkgs = pkgs: pkgs.rustBuilder.makePackageSet {
		rustVersion = "1.75.0";
		packageFun = import ./Cargo.nix;
	};
	
	default = (rustPkgs pkgs).workspace.{{.Name}} {};
    }
    `
)

type rustApp struct {
	ProjectName string
	RustVersion string
}

// GenerateRustApp generates default flake
func GenerateRustApp(fl *hcl2nix.RustApp, wr io.Writer) error {
	var rustVersion string
	if fl.RustVersion == ""{
		rustVersion = "1.75.0"
	} else {
		rustVersion = fl.RustVersion
	}
	data := rustApp{
		ProjectName: fl.ProjectName,
		RustVersion: rustVersion,
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
