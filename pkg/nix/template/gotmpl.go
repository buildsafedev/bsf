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
	   name = "{{ .Name }}";
	   src = {{ .SourcePath }};  
	   {{ if .DoCheck }}{{ else }}doCheck = false;{{ end }}
	   {{ if .VendorHash }}
		vendorHash = "{{ .VendorHash  }}";
		{{ else }}
		vendorHash = lib.fakeHash;
		{{ end }}
	   meta = with lib; {
		 description = "";
	   };
	   {{ if gt (len .LdFlags) 0}}
		ldflags = [
			{{ range $value := .LdFlags }}"{{ $value }}" {{ end }}
		];
	   {{ end }}
	   {{ if gt (len .Tags) 0 }}
		tags = [
			{{ range $value := .Tags }}"{{ $value }}" {{ end }}
		];
	   {{ end }}
	 }
	`
)

type goModule struct {
	Name       string
	SourcePath string
	LdFlags    []string
	Tags       []string
	VendorHash template.HTML
	DoCheck    bool
	Meta       *hcl2nix.Meta
}

// GenerateGoModule generates default flake
func GenerateGoModule(fl *hcl2nix.GoModule, wr io.Writer) error {
	data := goModule{
		Name:       fl.Name,
		SourcePath: fl.SourcePath,
		DoCheck:    fl.DoCheck,

		// Convert VendorHash to HTML to avoid escaping
		VendorHash: template.HTML(fl.VendorHash),
	}

	if len(fl.LdFlags) != 0 {
		data.LdFlags = fl.LdFlags
	}
	if len(fl.Tags) != 0 {
		data.Tags = fl.Tags
	}

	t, err := template.New(golangTmpl).Parse(golangTmpl)
	if err != nil {
		return err
	}

	err = t.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}
