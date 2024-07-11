package init

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
)

func generatehcl2NixConf(pt langdetect.ProjectType, pd *langdetect.ProjectDetails) (hcl2nix.Config, error) {
	switch pt {
	case langdetect.GoModule:
		return genGoModuleConf(pd), nil
	case langdetect.PythonPoetry:
		return genPythonPoetryConf(), nil
	case langdetect.RustCargo:
		config, err := genRustCargoConf()
		if err != nil {
			return hcl2nix.Config{}, err
		}
		return config, nil
	case langdetect.JsNpm:
		config, err := genJsNpmConf()
		if err != nil {
			return hcl2nix.Config{}, err
		}
		return config, nil
	default:
		return generateEmptyConf(), nil
	}
}

func generateEmptyConf() hcl2nix.Config {
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: []string{""},
			Runtime:     []string{"cacert@3.95"},
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact: "pkgs",
				Name:     "bsf-base-image",
				Cmd:      []string{},
				Entrypoint: []string{},
				EnvVars:    []string{},
				ExposedPorts: []string{},
				ImportConfigs: []string{},
			},
		},
	}
}
func genRustCargoConf() (hcl2nix.Config, error) {
	content, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return hcl2nix.Config{}, fmt.Errorf("error reading file: %v", err)
	}

	packageNameRegex, err := regexp.Compile(`name = "(.*?)"`)
	if err != nil {
		return hcl2nix.Config{}, fmt.Errorf("error fetching project name: %v", err)
	}

	match := packageNameRegex.FindStringSubmatch(string(content))

	var CrateName string
	if len(match) >= 2 {
		CrateName = match[1]
	} else {
		CrateName = "my-project"
	}
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: []string{"cargo@1.75.0"},
			Runtime:     []string{"cacert@3.95"},
		},
		RustApp: &hcl2nix.RustApp{
			WorkspaceSrc: "./.",
			CrateName:    CrateName,
			RustVersion:  "1.75.0",
			Release:      true,
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact: "pkgs",
				Name:     "bsf-base-image",
				Cmd:      []string{},
				Entrypoint: []string{},
				EnvVars:    []string{},
				ExposedPorts: []string{},
				ImportConfigs: []string{},
			},
		},
	}, nil
}

func genPythonPoetryConf() hcl2nix.Config {
	// TODO: maybe we should note down the path of the poetry.lock file and use it here.
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: []string{"python3@3.12.2", "poetry@1.8.2"},
			Runtime:     []string{"cacert@3.95"},
		},
		PoetryApp: &hcl2nix.PoetryApp{
			ProjectDir:   "./.",
			Src:          "./.",
			Pyproject:    "./pyproject.toml",
			Poetrylock:   "./poetry.lock",
			PreferWheels: false,
			CheckGroups:  []string{"dev"},
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact: "pkgs",
				Name:     "bsf-base-image",
				Cmd:      []string{},
				Entrypoint: []string{},
				EnvVars:    []string{},
				ExposedPorts: []string{},
				ImportConfigs: []string{},
			},
		},
	}
}

func genGoModuleConf(pd *langdetect.ProjectDetails) hcl2nix.Config {
	var name, entrypoint string
	if pd != nil {
		name = pd.Name
		entrypoint = pd.Entrypoint
		if entrypoint == "" {
			entrypoint = "./."
		}

	}
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: []string{"go@1.22.3", "gotools@0.18.0", "delve@1.22.1"},
			// todo: maybe we should dynamically inject the latest version of such runtime packages(cacert)?
			Runtime: []string{"cacert@3.95"},
		},
		GoModule: &hcl2nix.GoModule{
			Name:       name,
			SourcePath: entrypoint,
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact: "pkgs",
				Name:     "bsf-base-image",
				Cmd:      []string{},
				Entrypoint: []string{},
				EnvVars:    []string{},
				ExposedPorts: []string{},
				ImportConfigs: []string{},
			},
		},
	}

}

func genJsNpmConf() (hcl2nix.Config, error) {
	data, err := os.ReadFile("package-lock.json")
	if err != nil {
		return hcl2nix.Config{}, fmt.Errorf("error reading file: %v", err)
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return hcl2nix.Config{}, fmt.Errorf("error parsing json data: %v", err)
	}

	name, ok := jsonData["name"].(string)
	if !ok {
		return hcl2nix.Config{}, fmt.Errorf("error fetching project name: %v", err)
	}
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: []string{"nodejs@20.11.1"},
			Runtime:     []string{"cacert@3.95"},
		},
		JsNpmApp: &hcl2nix.JsNpmApp{
			PackageName: name,
			PackageRoot: "./.",
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact: "pkgs",
				Name:     "bsf-base-image",
				Cmd:      []string{},
				Entrypoint: []string{},
				EnvVars:    []string{},
				ExposedPorts: []string{},
				ImportConfigs: []string{},
			},
		},
	}, nil
}
