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
	Base          bool
}

const (
    ociTmpl = `
	ociImages = forEachSupportedSystem ({ pkgs, nix2containerPkgs, system , ...}: {
		{{range $artifact := .}}
		{{ if ne ($artifact.Base) true }}
		ociImage_{{$artifact.Artifact}}_app = nix2containerPkgs.nix2container.buildImage {
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
			layers = builtins.concatLists [
				(map (pkg: nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ pkg ];
				}) inputs.self.runtimeEnvsForOCI.${system}.runtime)

				{{range $config := $artifact.ImportConfigs}}
				(nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ inputs.self.configs.${system}.config_{{ . }} ];
				}),
				{{end}}

			];     
		};

		ociImage_{{$artifact.Artifact}}_app_with_dev = nix2containerPkgs.nix2container.buildImage {
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
			layers = builtins.concatLists [
				(map (pkg: nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ pkg ];
				}) inputs.self.runtimeEnvsForOCI.${system}.runtime)

				{{range $config := $artifact.ImportConfigs}}
				(nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ inputs.self.configs.${system}.config_{{ . }} ];
				}),
				{{end}}

				(map (pkg: nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ pkg ];
				}) inputs.self.devEnvsForOCI.${system}.development)
			];     
		};
		{{end}}

		{{ if ($artifact.Base)}}
		ociImage_{{$artifact.Artifact}}_runtime = nix2containerPkgs.nix2container.buildImage {
			name = "{{$artifact.Name}}";
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
			layers = builtins.concatLists [
				(map (pkg: nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ pkg ];
				}) inputs.self.runtimeEnvsForOCI.${system}.runtime)

				{{range $config := $artifact.ImportConfigs}}
				(nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ inputs.self.configs.${system}.config_{{ . }} ];
				}),
				{{end}}

			];
		};

		ociImage_{{$artifact.Artifact}}_dev = nix2containerPkgs.nix2container.buildImage {
			name = "{{$artifact.Name}}";
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
			layers = builtins.concatLists [
				(map (pkg: nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ pkg ];
				}) inputs.self.runtimeEnvsForOCI.${system}.runtime)

				{{range $config := $artifact.ImportConfigs}}
				(nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ inputs.self.configs.${system}.config_{{ . }} ];
				}),
				{{end}}

				(map (pkg: nix2containerPkgs.nix2container.buildLayer {
					copyToRoot = [ pkg ];
				}) inputs.self.devEnvsForOCI.${system}.development)
			];
		};
		{{end}}
		{{end}}

		{{range $artifact := .}}
		{{ if ne ($artifact.Base) true }}
		ociImage_{{$artifact.Artifact}}_app-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImages.${system}.ociImage_{{$artifact.Artifact}}_app.copyTo}/bin/copy-to dir:$out";
		ociImage_{{$artifact.Artifact}}_app_with_dev-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImages.${system}.ociImage_{{$artifact.Artifact}}_app_with_dev.copyTo}/bin/copy-to dir:$out";
		{{end}}
		{{ if ($artifact.Base)}}
		ociImage_{{$artifact.Artifact}}_runtime-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImages.${system}.ociImage_{{$artifact.Artifact}}_runtime.copyTo}/bin/copy-to dir:$out";
		ociImage_{{$artifact.Artifact}}_dev-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImages.${system}.ociImage_{{$artifact.Artifact}}_dev.copyTo}/bin/copy-to dir:$out";
		{{end}}
		{{end}}
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
		}
		if ociArtifact.IsBase {
			converted[i].Base = true
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