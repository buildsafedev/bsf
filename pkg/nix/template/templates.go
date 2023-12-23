package template

import (
	"html/template"
	"os"
)

// Flake holds flake parameters
type Flake struct {
	Description      string
	PackageInputUrls []string
	PackageInputs    map[string]string
	DevPackages      map[string]string
	RuntimePackages  map[string]string
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
		{{range .PackageInputUrls}} nixpkgs-{{ .}}.url = "github.com/nixos/nixpkgs/{{ . }}";
		{{ end }}	
	};
	
	outputs = { self, {{range .PackageInputUrls}} nixpkgs-{{ .}}, 
	{{end}} }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		{{ range .PackageInputUrls }} nixpkgs-{{ .}}-pkgs = import nixpkgs-{{ .}} { inherit system; };
		{{ end }}
	  });
	in {
	  packages = forEachSupportedSystem ({ {{ range .PackageInputUrls }} nixpkgs-{{ .}}-pkgs, 
		{{ end }} }: {
		default = pkgs.callPackage ./default.nix {
		  {{range $key, $value :=.PackageInputs}} {{ $key }} = {{ $value }};{{ end }}
		};
	  });
	
	  devShells = forEachSupportedSystem ({ {{ range .PackageInputUrls }} nixpkgs-{{ .}}-pkgs, 
		{{ end }} }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			{{ range $key, $value :=.DevPackages }}nixpkgs-{{ $value  }}-pkgs.{{ $key }}  
			{{ end }}
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ {{ range .PackageInputUrls }} nixpkgs-{{ .}}-pkgs, {{ end }} }: {
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

func templateMain(fl Flake) (*string, error) {
	t, err := template.New("main").Parse(mainTmpl)
	if err != nil {
		return nil, err
	}

	err = t.Execute(os.Stdout, fl)
	if err != nil {
		return nil, err
	}

	return &t.Tree.Name, nil
}
