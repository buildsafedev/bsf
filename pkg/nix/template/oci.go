package template

import (
	"bytes"
	"text/template"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

// OCIArtifact holds parameters for OCI artifacts
type OCIArtifact struct {
	Artifact      string
	Name          string
	Cmd           []string
	Entrypoint    []string
	EnvVars       []string
	ImportConfigs []string
	ExposedPorts  []string
	RuntimeDeps   bool
	DevDeps       bool
	Base          bool
}

const (
	ociTmpl = `
	ociImages = forEachSupportedSystem ({ pkgs, nix2containerPkgs, system , ...}: {
		{{range $artifact := .}}
		ociImage_{{$artifact.Artifact}} =  nix2containerPkgs.nix2container.buildImage {
		name = "{{$artifact.Name}}";
		{{if ne .Base true}}
		copyToRoot = [ inputs.self.packages.${system}.default ];
		{{end}}
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
				{{ if (.RuntimeDeps)}}
			    inputs.self.runtimeEnvs.${system}.runtime
				{{end}}
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
	ociImage_{{$artifact.Artifact}}-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImages.${system}.ociImage_{{$artifact.Artifact}}.copyTo}/bin/copy-to dir:$out";{{end}}
	});
	`
)

func hclOCIToOCIArtifact(ociArtifacts []hcl2nix.OCIArtifact) []OCIArtifact {
	converted := make([]OCIArtifact, len(ociArtifacts))

	for i, ociArtifact := range ociArtifacts {
		converted[i] = OCIArtifact{
			Artifact:      ociArtifact.Artifact,
			Name:          ociArtifact.Name,
			Cmd:           ociArtifact.Cmd,
			Entrypoint:    ociArtifact.Entrypoint,
			EnvVars:       ociArtifact.EnvVars,
			ImportConfigs: ociArtifact.ImportConfigs,
			ExposedPorts:  ociArtifact.ExposedPorts,
			DevDeps:       ociArtifact.DevDeps,
		}
		if ociArtifact.Artifact == "pkgs" {
			converted[i].Base = true
			if ociArtifact.DevDeps {
				converted[i].RuntimeDeps = false
			} else {
				converted[i].RuntimeDeps = true
			}
		} else {
			converted[i].RuntimeDeps = true
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
