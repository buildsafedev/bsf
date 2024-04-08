package template

import (
	"html/template"
	"io"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

// Flake holds flake parameters
type Flake struct {
	Description         string
	Language            string
	NixPackageRevisions []string
	DevPackages         map[string]string
	RuntimePackages     map[string]string
	RustArguments       RustApp
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
		
		{{if eq .Language "PythonPoetry"}} poetry2nix = {
			url = "github:nix-community/poetry2nix";
			inputs.nixpkgs.follows = "nixpkgs";
		  }; {{end}}
		
		{{if eq .Language "RustCargo"}}
		 cargo2nix.url = "github:cargo2nix/cargo2nix/release-0.11.0";
    	 nixpkgs.follows = "cargo2nix/nixpkgs";{{end}}
	};
	
	outputs = { self, nixpkgs, 
	{{if eq .Language "GoModule"}} gomod2nix, {{end}}
	{{ if eq .Language "PythonPoetry"}} poetry2nix, {{end}}
	{{ if eq .Language "RustCargo"}} cargo2nix, {{end}}
	{{range .NixPackageRevisions}} nixpkgs-{{ .}}, 
	{{end}} }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  {{ if eq .Language "RustCargo"}}
	  rustPkgs = pkgs: pkgs.rustBuilder.makePackageSet {
		packageFun = import ./Cargo.nix;
		workspaceSrc = ../.;
		{{ if ne .RustArguments.RustVersion ""}}
		rustVersion = "{{ .RustArguments.RustVersion }}"; {{ end }}
		{{ if ne .RustArguments.RustToolChain ""}}
		rustToolchain = "{{ .RustArguments.RustToolChain }}"; {{ end }}
		{{ if ne .RustArguments.RustChannel ""}}
		rustChannel = "{{ .RustArguments.RustChannel }}"; {{ end }}
		{{ if ne .RustArguments.RustProfile ""}}
		rustProfile = "{{ .RustArguments.RustProfile }}"; {{ end }}
		{{ if gt (len .RustArguments.ExtraRustComponents) 0}}
		extraRustComponenets = [{{ range $value := .RustArguments.ExtraRustComponents }}"{{ $value }}",{{ end }}];{{ end }}
		{{ if ne .RustArguments.Release true}}
		release = {{ .RustArguments.Release }}; {{ end }}
	  }; {{end}}
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		{{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs = import nixpkgs-{{ .}} { inherit system; };
		{{ end }}
		{{if eq .Language "GoModule"}} buildGoApplication = gomod2nix.legacyPackages.${system}.buildGoApplication;{{end}}
		pkgs = import nixpkgs { inherit system; {{ if eq .Language "RustCargo"}} overlays = [cargo2nix.overlays.default]; {{end}} };
		{{if eq .Language "PythonPoetry"}} inherit (poetry2nix.lib.mkPoetry2Nix { pkgs = nixpkgs.legacyPackages.${system}; }) mkPoetryApplication; {{end}}
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,
		{{if eq .Language "GoModule"}} buildGoApplication, {{end}}
		{{if eq .Language "PythonPoetry"}} mkPoetryApplication, {{end}}
		{{ range .NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, 
		{{ end }} }: {
		default = pkgs.callPackage ./default.nix {
			{{if eq .Language "GoModule"}} inherit buildGoApplication;
			go = pkgs.go_1_22; {{end}}
			{{if eq .Language "PythonPoetry"}} inherit mkPoetryApplication; {{end}}
			{{if eq .Language "RustCargo"}}
			 inherit pkgs;
             inherit rustPkgs;
			{{end}}
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs, 
		{{if eq .Language "GoModule"}} buildGoApplication, {{end}}
		{{if eq .Language "PythonPoetry"}} mkPoetryApplication, {{end}}
		{{if eq .Language "RustCargo"}} rustPkgs, {{end}}
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
		{{if eq .Language "PythonPoetry"}} mkPoetryApplication, {{end}}
		{{if eq .Language "RustCargo"}} rustPkgs, {{end}}
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
		{{if eq .Language "PythonPoetry"}} mkPoetryApplication, {{end}}
		{{if eq .Language "RustCargo"}} rustPkgs, {{end}}
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
func GenerateFlake(fl Flake, wr io.Writer, conf *hcl2nix.Config) error {
	if conf.RustApp != nil {
		fl.RustArguments = RustApp{
			RustVersion:         conf.RustApp.RustVersion,
			RustToolChain:       conf.RustApp.RustToolChain,
			RustChannel:         conf.RustApp.RustChannel,
			RustProfile:         conf.RustApp.RustProfile,
			ExtraRustComponents: conf.RustApp.ExtraRustComponents,
			Release:             conf.RustApp.Release,
		}
	}

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
