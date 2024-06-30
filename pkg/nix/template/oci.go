package template

import (
	"bytes"
	"text/template"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

// OCIArtifact holds parameters for OCI artifacts
type OCIArtifact struct {
	Environment   string
	Name          string
	Cmd           []string
	Entrypoint    []string
	EnvVars       []string
	ImportConfigs []string
	ExposedPorts  []string
	DevDeps       bool
}

// OCIArtifactforBase holds parameters for base OCI Artifacts
type OCIArtifactforBase struct {
	Name    string
	Runtime bool
	Dev     bool
}

const (
	ociTmpl = `
	ociImages = forEachSupportedSystem ({ pkgs, nix2containerPkgs, system , ...}: {
		{{range $artifact := .}}
		ociImage_{{$artifact.Environment}} =  nix2containerPkgs.nix2container.buildImage {
		name = "{{$artifact.Name}}";
		copyToRoot = [ inputs.self.packages.${system}.default ];
		  config = {

			cmd = [ {{range $c := $artifact.Cmd}}
				"{{.}}" {{end}} ];

			entrypoint = [ {{range $c := $artifact.Entrypoint}}
				"{{.}}" {{end}} ];

		  env = [
			{{range $env := $artifact.EnvVars}}
		     "{{ . }}"{{end}}
		  ];
		  ExposedPorts = {
			{{ range $port := $artifact.ExposedPorts}}
			 "{{ . }}"={}; {{end}}
		  };
		  };
		 maxLayers = 100;
		 layers = [
			 (nix2containerPkgs.nix2container.buildLayer { 
				copyToRoot = [
			    inputs.self.runtimeEnvs.${system}.runtime
			    {{range $config := $artifact.ImportConfigs}}
				inputs.self.configs.${system}.config_{{ . }} {{end}}
				{{ if (.DevDeps)}}
				inputs.self.devEnvs.${system}.development
				{{end}}
			  ];
			 })
		  ];      
	};
	{{end}}
	{{range $artifact := .}}
	ociImage_{{$artifact.Environment}}-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImages.${system}.ociImage_{{$artifact.Environment}}.copyTo}/bin/copy-to dir:$out";{{end}}
	});
	`
)

const (
	ociTmplForBase = `
	ociImages = forEachSupportedSystem ({ pkgs, nix2containerPkgs, system , ...}: {
		ociImage_{{.Name}} = nix2containerPkgs.nix2container.buildImage {
		name = "{{.Name}}";
		 maxLayers = 100;
		 layers = [
			 (nix2containerPkgs.nix2container.buildLayer { 
				copyToRoot = [
				{{if .Runtime}}
			    inputs.self.runtimeEnvs.${system}.runtime
				{{end}}
				{{if .Dev}}
				inputs.self.devEnvs.${system}.development
				{{end}}
			  ];
			 })
		  ];      
	};
	ociImage_{{.Name}}-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImages.${system}.ociImage_{{.Name}}.copyTo}/bin/copy-to dir:$out";
	});
	`
)

func hclOCIToOCIArtifact(ociArtifacts []hcl2nix.OCIArtifact) []OCIArtifact {
	converted := make([]OCIArtifact, len(ociArtifacts))

	for i, ociArtifact := range ociArtifacts {
		converted[i] = OCIArtifact{
			Environment:   ociArtifact.Environment,
			Name:          ociArtifact.Name,
			Cmd:           ociArtifact.Cmd,
			Entrypoint:    ociArtifact.Entrypoint,
			EnvVars:       ociArtifact.EnvVars,
			ImportConfigs: ociArtifact.ImportConfigs,
			ExposedPorts:  ociArtifact.ExposedPorts,
			DevDeps:       ociArtifact.DevDeps,
		}
	}

	return converted

}

// GenerateOCIAttr generates the Nix attribute set for oci artifacts
func GenerateOCIAttr(artifacts []OCIArtifact) (*string, error) {
	tmpl, err := template.New("ociAttr").Funcs(template.FuncMap{
		"quote": quote,
	}).
		Parse(ociTmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, artifacts)
	if err != nil {
		return nil, err
	}

	result := buf.String()
	return &result, nil
}

// GenerateOCIAttrForBase generates the Nix attribute set for oci artifacts for base image
func GenerateOCIAttrForBase(artifact OCIArtifactforBase) (*string, error) {
	tmpl, err := template.New("ociAttrForBase").Funcs(template.FuncMap{
		"quote": quote,
	}).
		Parse(ociTmplForBase)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, artifact)
	if err != nil {
		return nil, err
	}

	result := buf.String()
	return &result, nil
}
