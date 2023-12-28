package template

import (
	"html/template"
	"io"
)

// Flake holds flake parameters
type Flake struct {
	Description         string
	NixPackageRevisions []string
	PackageInputs       map[string]string
	DevPackages         map[string]string
	RuntimePackages     map[string]string
}

// Meta holds meta parameters
type Meta struct {
	Description string
}

const (
	mainTmpl = `
{
	description = "{{.Description }}";
	
	inputs = {
		{{range .NixPackageRevisions}} nixpkgs-{{ .}}.url = "github.com/nixos/nixpkgs/{{ . }}";
		{{ end }}	
	};
	
	outputs = { self, {{range .NixPackageRevisions}} nixpkgs-{{ .}}, 
	{{end}} }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		{{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs = import nixpkgs-{{ .}} { inherit system; };
		{{ end }}
	  });
	in {
	  packages = forEachSupportedSystem ({ {{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, 
		{{ end }} }: {
		default = pkgs.callPackage ./default.nix {
		  {{range $key, $value :=.PackageInputs}} {{ $key }} = {{ $value }};{{ end }}
		};
	  });
	
	  devShells = forEachSupportedSystem ({ {{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, 
		{{ end }} }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			{{ range $key, $value :=.DevPackages }}nixpkgs-{{ $value  }}-pkgs.{{ $key }}  
			{{ end }}
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ {{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, {{ end }} }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			{{ range $key, $value := .RuntimePackages }}nixpkgs-{{ $value  }}-pkgs.{{$key}}   
			{{ end }}
		   ];
		};
	});
	};
}
`
)

// GenerateDefaultFlake generates default flake
func GenerateDefaultFlake(fl Flake, wr io.Writer) error {
	t, err := template.New("main").Parse(mainTmpl)
	if err != nil {
		return err
	}

	err = t.Execute(wr, fl)
	if err != nil {
		return err
	}

	return nil
}
