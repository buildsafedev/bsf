package template

import (
	"html/template"
	"io"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

const (
	golangTmpl = `
	{ pkgs ? (
		let
		  inherit (builtins) fetchTree fromJSON readFile;
		  inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
		in
		import (fetchTree nixpkgs.locked) {
		  overlays = [
			(import "${fetchTree gomod2nix.locked}/overlay.nix")
		  ];
		}
	  )
	, buildGoApplication ? pkgs.buildGoApplication
	}:
	
	buildGoApplication {
	  pname = "{{ .Name }}";
	  version = "0.1";
	  pwd = ./.;
	  src = {{ .SourcePath }};  
	  modules = ./gomod2nix.toml;
	  {{ if .DoCheck }}{{ else }}doCheck = false;{{ end }}
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
	DoCheck    bool
}

// GenerateGoModule generates default flake
func GenerateGoModule(fl *hcl2nix.GoModule, wr io.Writer) error {
	data := goModule{
		Name:       fl.Name,
		SourcePath: fl.SourcePath,
		DoCheck:    fl.DoCheck,
	}

	if len(fl.LdFlags) != 0 {
		data.LdFlags = fl.LdFlags
	}
	if len(fl.Tags) != 0 {
		data.Tags = fl.Tags
	}

	t, err := template.New("go").Parse(golangTmpl)
	if err != nil {
		return err
	}

	err = t.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}
