package template

import (
	"html/template"
	"io"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

const (
	golangTmpl = `
	{
	   lib,
	   stdenv,
	   buildGoModule,
	   ... 
	 }: buildGoModule {
	   name = "";
	   src = ../.;  
	   {{ if .VendorHash }}
		vendorSha256 = "{{ .VendorHash }}";
		{{ else }}
		vendorHash = lib.fakeHash;
		{{ end }}
	   meta = with lib; {
		 description = "";
	   };
	 }
	`
)

// GenerateGoModule generates default flake
func GenerateGoModule(fl *hcl2nix.GoModule, wr io.Writer) error {
	t, err := template.New("gomodule").Parse(golangTmpl)
	if err != nil {
		return err
	}

	err = t.Execute(wr, fl)
	if err != nil {
		return err
	}

	return nil
}
