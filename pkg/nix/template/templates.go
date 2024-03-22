package template

import (
	"html/template"
	"io"
)

// Flake holds flake parameters
type Flake struct {
	Description         string
	Language            string
	NixPackageRevisions []string
	DevPackages         map[string]string
	RuntimePackages     map[string]string
}

// todo: maybe we could let power users inject their own templates
const (
	mainTmpl = `
{
	description = "{{.Description }}";
	
	inputs = {
		{{range .NixPackageRevisions}} nixpkgs-{{ .}}.url = "github:nixos/nixpkgs/{{ . }}";
		{{ end }}	
		nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
		{{if eq .Language "GoModule"}} gomod2nix.url = "github:nix-community/gomod2nix";
		gomod2nix.inputs.nixpkgs.follows = "nixpkgs";{{end}}
	};
	
	outputs = { self, nixpkgs, 
	{{if eq .Language "GoModule"}} gomod2nix, {{end}}
	{{range .NixPackageRevisions}} nixpkgs-{{ .}}, 
	{{end}} }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		{{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs = import nixpkgs-{{ .}} { inherit system; };
		{{ end }}
		{{if eq .Language "GoModule"}} buildGoApplication = gomod2nix.legacyPackages.${system}.buildGoApplication;{{end}}
		pkgs = import nixpkgs { inherit system; };
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,
		{{if eq .Language "GoModule"}} buildGoApplication, {{end}}
		{{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, 
		{{ end }} }: {
		default = pkgs.callPackage ./default.nix {
			{{if eq .Language "GoModule"}} inherit buildGoApplication;
			go = pkgs.go_1_22; {{end}}
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs, 
		{{if eq .Language "GoModule"}} buildGoApplication, {{end}}
		{{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, 
		{{ end }} }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			{{ range $key, $value :=.DevPackages }}nixpkgs-{{ $value  }}-pkgs.{{ $key }}  
			{{ end }}
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ pkgs,
		{{if eq .Language "GoModule"}} buildGoApplication, {{end}}
		{{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, {{ end }} }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			{{ range $key, $value := .RuntimePackages }}nixpkgs-{{ $value  }}-pkgs.{{$key}}   
			{{ end }}
		   ];
		};
	   });

	   devEnvs = forEachSupportedSystem ({ pkgs,
		{{if eq .Language "GoModule"}} buildGoApplication, {{end}}
	   {{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, {{ end }} }: {
		development = pkgs.buildEnv {
		  name = "devenv";
		  paths = [ 
			{{ range $key, $value :=.DevPackages }}nixpkgs-{{ $value  }}-pkgs.{{ $key }}  
			{{ end }}
		   ];
		};
	   });
	};
}
`
)

// GenerateFlake generates default flake
func GenerateFlake(fl Flake, wr io.Writer) error {
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

// parentFolder returns the parent folder of the given path. ex: ( ./ -> ../ )
func parentFolder(s string) string {
	return "." + s
}
