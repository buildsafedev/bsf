package template

import (
	"bytes"
	"text/template"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

// OCIArtifact holds parameters for OCI artifacts
type OCIArtifact struct {
	Artifact            string
	Name                string
	Layers              []string
	Cmd                 []string
	Entrypoint          []string
	EnvVars             []string
	ImportConfigs       []string
	ExposedPorts        []string
	Base                bool
	NixPackageRevisions []string
}

const (
	ociTmpl = `
{{range $artifact := .}}
ociImage_{{$artifact.Artifact}} = forEachSupportedSystem ({ pkgs, nix2containerPkgs, system , {{ range $artifact.NixPackageRevisions }} nixpkgs-{{ .}}-pkgs, {{ end }} ... }: {
  {{ if ne ($artifact.Base) true }}
  ociImage_{{$artifact.Artifact}}_app = nix2containerPkgs.nix2container.buildImage {
    name = "{{$artifact.Name}}";
    copyToRoot = [ inputs.self.packages.${system}.default ];
    config = {
      cmd = [ {{range $c := $artifact.Cmd}} "{{.}}" {{end}} ];
      entrypoint = [ {{range $c := $artifact.Entrypoint}} "{{.}}" {{end}} ];
      env = [
        {{range $env := $artifact.EnvVars}} "{{ . }}"{{end}}
      ];
      ExposedPorts = {
        {{ range $port := $artifact.ExposedPorts}} "{{ . }}" = {}; {{end}}
      };
    };
    maxLayers = 100;
    layers = [
      {{range $layer := .Layers}} {{$layer}} {{end}}
      {{range $config := $artifact.ImportConfigs}}
      (nix2containerPkgs.nix2container.buildLayer {
        copyToRoot = [ inputs.self.configs.${system}.config_{{ . }} ];
      }),
      {{end}}
    ];
  };
  {{end}}

  {{ if ($artifact.Base)}}
  ociImage_{{$artifact.Artifact}}_base = nix2containerPkgs.nix2container.buildImage {
    name = "{{$artifact.Name}}";
    config = {
      cmd = [ {{range $c := $artifact.Cmd}} "{{.}}" {{end}} ];
      entrypoint = [ {{range $c := $artifact.Entrypoint}} "{{.}}" {{end}} ];
      env = [
        {{range $env := $artifact.EnvVars}} "{{ . }}"{{end}}
      ];
      ExposedPorts = {
        {{ range $port := $artifact.ExposedPorts}} "{{ . }}" = {}; {{end}}
      };
    };
    maxLayers = 100;
    layers = [
      {{range $layer := .Layers}} {{$layer}} {{end}}
      {{range $config := $artifact.ImportConfigs}}
      (nix2containerPkgs.nix2container.buildLayer {
        copyToRoot = [ inputs.self.configs.${system}.config_{{ . }} ];
      }),
      {{end}}
    ];
  };
  {{end}}

  {{ if ne ($artifact.Base) true }}
  ociImage_{{$artifact.Artifact}}_app-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImage.${system}.ociImage_{{$artifact.Artifact}}_app.copyTo}/bin/copy-to dir:$out";
  {{end}}
  {{ if ($artifact.Base)}}
  ociImage_{{$artifact.Artifact}}_base-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImage_{{$artifact.Artifact}}.${system}.ociImage_{{$artifact.Artifact}}_base.copyTo}/bin/copy-to dir:$out";
  {{end}}
  });
{{end}}

`

)

func hclOCIToOCIArtifact(ociArtifacts []hcl2nix.OCIArtifact, fl Flake) []OCIArtifact {
	converted := make([]OCIArtifact, len(ociArtifacts))

	for i, ociArtifact := range ociArtifacts {
		converted[i] = OCIArtifact{
			Artifact:            ociArtifact.Artifact,
			Name:                ociArtifact.Name,
			Layers:              getLayers(ociArtifact.Layers, fl),
			Cmd:                 ociArtifact.Cmd,
			Entrypoint:          ociArtifact.Entrypoint,
			EnvVars:             ociArtifact.EnvVars,
			ImportConfigs:       ociArtifact.ImportConfigs,
			ExposedPorts:        ociArtifact.ExposedPorts,
			NixPackageRevisions: fl.NixPackageRevisions,
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

func getLayers(layers []string, fl Flake) []string {
	reqPkgs := getReqPkgs(layers, fl)

	var layerBlocks []string

	for _, pkg := range reqPkgs {
		var layerBlock string

		if pkg == "pkgs.runtime" {
			layerBlock = `
			(nix2containerPkgs.nix2container.buildLayer { 
				copyToRoot = [
					inputs.self.runtimeEnvs.${system}.runtime
				];
			})`
		} else if pkg == "pkgs.dev" {
			layerBlock = `
			(nix2containerPkgs.nix2container.buildLayer { 
				copyToRoot = [
					inputs.self.devEnvs.${system}.development
				];
			})`
		} else {
			layerBlock = `
			(nix2containerPkgs.nix2container.buildLayer { 
				copyToRoot = [
					` + pkg + `
				];
			})`
		}

		layerBlocks = append(layerBlocks, layerBlock)
	}
	return layerBlocks
}

func getReqPkgs(layers []string, fl Flake) []string {
	var newLayers []string

	for _, l := range layers {
		if l == "split(pkgs.runtime)" {
			for key, value := range fl.RuntimePackages {
				newLayer := "nixpkgs-" + value + "-pkgs." + key
				newLayers = append(newLayers, newLayer)
			}
		} else if l == "split(pkgs.dev)" {
			for key, value := range fl.DevPackages {
				newLayer := "nixpkgs-" + value + "-pkgs." + key
				newLayers = append(newLayers, newLayer)
			}
		} else if len(l) > 8 && l[:9] == "pkgs.dev." {
			pkgName := l[9:]
			for key, value := range fl.DevPackages {
				if key == pkgName {
					newLayer := "nixpkgs-" + value + "-pkgs." + key
					newLayers = append(newLayers, newLayer)
					break
				}
			}
		} else if len(l) > 12 && l[:13] == "pkgs.runtime." {
			pkgName := l[13:]
			for key, value := range fl.RuntimePackages {
				if key == pkgName {
					newLayer := "nixpkgs-" + value + "-pkgs." + key
					newLayers = append(newLayers, newLayer)
					break
				}
			}
		} else {
			newLayers = append(newLayers, l)
		}
	}

	return newLayers
}
