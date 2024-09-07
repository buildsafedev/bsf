package template

import (
	"bytes"
	"strings"
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
  ociImage_{{$artifact.Artifact}}_app-as-dir = pkgs.runCommand "image-as-dir" { } "${inputs.self.ociImage_{{$artifact.Artifact}}.${system}.ociImage_{{$artifact.Artifact}}_app.copyTo}/bin/copy-to dir:$out";
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

	for _, pkgSet := range reqPkgs {
		layerBlock :=
			`(nix2containerPkgs.nix2container.buildLayer { 
			copyToRoot = [
				` + strings.Join(pkgSet, "\n") + `
			];
		})`
		layerBlocks = append(layerBlocks, layerBlock)
	}

	return layerBlocks
}

func getReqPkgs(layers []string, fl Flake) [][]string {
	var newLayers [][]string

	for _, l := range layers {
		if strings.Contains(l, " + ") {
			combinedLayer := handleCombinedLayers(l, fl)
			newLayers = append(newLayers, combinedLayer)
		} else {
			individualLayers := handleIndividualLayers(l, fl)
			newLayers = append(newLayers, individualLayers...)
		}
	}

	return newLayers
}

func handleCombinedLayers(layer string, fl Flake) []string {
	var combinedLayer []string
	parts := strings.Split(layer, " + ")

	for _, part := range parts {
		switch {
		case part == "split(packages.runtime)":
			for key, value := range fl.RuntimePackages {
				combinedLayer = append(combinedLayer, "nixpkgs-"+value+"-pkgs."+key)
			}
		case part == "split(packages.dev)":
			for key, value := range fl.DevPackages {
				combinedLayer = append(combinedLayer, "nixpkgs-"+value+"-pkgs."+key)
			}
		case part == "packages.runtime":
			combinedLayer = append(combinedLayer, "inputs.self.runtimeEnvs.${system}.runtime")
		case part == "packages.dev":
			combinedLayer = append(combinedLayer, "inputs.self.devEnvs.${system}.development")
		case part == "gomodule" || part == "rustapp" || part == "jsnpmapp" || part == "poetryapp":
			combinedLayer = append(combinedLayer, "inputs.self.packages.${system}.default")
		case strings.HasPrefix(part, "packages.dev."):
			pkgName := strings.TrimPrefix(part, "packages.dev.")
			if value, exists := fl.DevPackages[pkgName]; exists {
				combinedLayer = append(combinedLayer, "nixpkgs-"+value+"-pkgs."+pkgName)
			}
		case strings.HasPrefix(part, "packages.runtime."):
			pkgName := strings.TrimPrefix(part, "packages.runtime.")
			if value, exists := fl.RuntimePackages[pkgName]; exists {
				combinedLayer = append(combinedLayer, "nixpkgs-"+value+"-pkgs."+pkgName)
			}
		}
	}

	return combinedLayer
}

func handleIndividualLayers(layer string, fl Flake) [][]string {
	var newLayers [][]string

	switch {
	case layer == "split(packages.runtime)":
		for key, value := range fl.RuntimePackages {
			newLayers = append(newLayers, []string{"nixpkgs-" + value + "-pkgs." + key})
		}
	case layer == "split(packages.dev)":
		for key, value := range fl.DevPackages {
			newLayers = append(newLayers, []string{"nixpkgs-" + value + "-pkgs." + key})
		}
	case layer == "packages.runtime":
		newLayers = append(newLayers, []string{"inputs.self.runtimeEnvs.${system}.runtime"})
	case layer == "packages.dev":
	case layer == "gomodule" || layer == "rustapp" || layer == "jsnpmapp" || layer == "poetryapp":
		newLayers = append(newLayers, []string{"inputs.self.packages.${system}.default"})
	case strings.HasPrefix(layer, "packages.dev."):
		pkgName := strings.TrimPrefix(layer, "packages.dev.")
		if value, exists := fl.DevPackages[pkgName]; exists {
			newLayers = append(newLayers, []string{"nixpkgs-" + value + "-pkgs." + pkgName})
		}
	case strings.HasPrefix(layer, "packages.runtime."):
		pkgName := strings.TrimPrefix(layer, "packages.runtime.")
		if value, exists := fl.RuntimePackages[pkgName]; exists {
			newLayers = append(newLayers, []string{"nixpkgs-" + value + "-pkgs." + pkgName})
		}
	default:
		newLayers = append(newLayers, []string{layer})
	}

	return newLayers
}
