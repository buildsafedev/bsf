package init

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
)

var (
	commonDevDeps = []string{"coreutils-full@9.5", "bash@5.2.15"}
	commonRTDeps  = []string{"cacert@3.95"}
	baseImageName = "ttl.sh/base"
)

func generatehcl2NixConf(pt langdetect.ProjectType, pd *langdetect.ProjectDetails) (hcl2nix.Config, error) {
	fmt.Println("IN GENERATE HCL NIX CONFIG -------------------------------------------")
	fmt.Println("NAME: ", pd.Name)
	switch pt {
	case langdetect.GoModule:
		return genGoModuleConf(pd), nil
	case langdetect.PythonPoetry:
		return genPythonPoetryConf(pd), nil
	case langdetect.RustCargo:
		config, err := genRustCargoConf(pd)
		if err != nil {
			return hcl2nix.Config{}, err
		}
		return config, nil
	case langdetect.JsNpm:
		config, err := genJsNpmConf(pd)
		if err != nil {
			return hcl2nix.Config{}, err
		}
		return config, nil
	default:
		return generateEmptyConf(pd), nil
	}
}

func generateEmptyConf(pd *langdetect.ProjectDetails) hcl2nix.Config {
	if pd.Name == "" {
		pd.Name = "expl"
	}
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: commonDevDeps,
			Runtime:     commonRTDeps,
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact:      "pkgs",
				Name:          baseImageName,
				Cmd:           []string{},
				Entrypoint:    []string{},
				EnvVars:       []string{},
				ExposedPorts:  []string{},
				ImportConfigs: []string{},
			},
		},
	}
}
func genRustCargoConf(pd *langdetect.ProjectDetails) (hcl2nix.Config, error) {
	if pd.Name == "" {
		pd.Name = "expl"
	}
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

	rustDevDeps := append(commonDevDeps, "cargo@1.75.0")
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: rustDevDeps,
			Runtime:     commonRTDeps,
		},
		RustApp: &hcl2nix.RustApp{
			WorkspaceSrc: "./.",
			CrateName:    CrateName,
			RustVersion:  "1.75.0",
			Release:      true,
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact:      "pkgs",
				Name:          baseImageName,
				Cmd:           []string{},
				Entrypoint:    []string{},
				EnvVars:       []string{},
				ExposedPorts:  []string{},
				ImportConfigs: []string{},
			},
		},
	}, nil
}

func genPythonPoetryConf(pd *langdetect.ProjectDetails) hcl2nix.Config {
	// TODO: maybe we should note down the path of the poetry.lock file and use it here.
	poetryDevDeps := append(commonDevDeps, "python3@3.12.2", "poetry@1.8.2")
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: poetryDevDeps,
			Runtime:     commonRTDeps,
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
				Artifact:      "pkgs",

				Name:          baseImageName,
				Cmd:           []string{},
				Entrypoint:    []string{},
				EnvVars:       []string{},
				ExposedPorts:  []string{},
				ImportConfigs: []string{},
			},
		},
	}
}

func genGoModuleConf(pd *langdetect.ProjectDetails) hcl2nix.Config {
	if pd.Name == "" {
		pd.Name = "expl"
	}
	var name, entrypoint string
	if pd != nil {
		name = pd.Name
		entrypoint = pd.Entrypoint
		if entrypoint == "" {
			entrypoint = "./."
		}

	}

	goDevDeps := append(commonDevDeps, "go@1.22.3", "gotools@0.18.0", "delve@1.22.1")
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: goDevDeps,
			// todo: maybe we should dynamically inject the latest version of such runtime packages(cacert)?
			Runtime: commonRTDeps,
		},
		GoModule: &hcl2nix.GoModule{
			Name:       name,
			SourcePath: entrypoint,
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact:      "pkgs",
				Name:          baseImageName,
				Cmd:           []string{},
				Entrypoint:    []string{},
				EnvVars:       []string{},
				ExposedPorts:  []string{},
				ImportConfigs: []string{},
			},
		},
	}

}

func genJsNpmConf(pd *langdetect.ProjectDetails) (hcl2nix.Config, error) {
	if pd.Name == "" {
		pd.Name = "expl"
	}
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

	nodeDevDeps := append(commonDevDeps, "nodejs@20.11.1")
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: nodeDevDeps,
			Runtime:     commonRTDeps,
		},
		JsNpmApp: &hcl2nix.JsNpmApp{
			PackageName: name,
			PackageRoot: "./.",
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact:      "pkgs",
				Name:          baseImageName,
				Cmd:           []string{},
				Entrypoint:    []string{},
				EnvVars:       []string{},
				ExposedPorts:  []string{},
				ImportConfigs: []string{},
			},
		},
	}, nil
}
