package init

import (
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
		return genPythonPoetryConf(pd), nil
	case langdetect.RustCargo:
		config, err := genRustCargoConf(pd)
		if err != nil {
			return hcl2nix.Config{}, err
		}
		return config, nil
	default:
		return hcl2nix.Config{
			Packages: hcl2nix.Packages{},
		}, nil
	}
}

func genRustCargoConf(pd *langdetect.ProjectDetails) (hcl2nix.Config, error) {
	content, err := os.ReadFile("Cargo.toml")
	if err != nil {
		return hcl2nix.Config{}, fmt.Errorf("Error reading file:", err)
	}

	packageNameRegex, err := regexp.Compile(`name = "(.*?)"`)
	if err != nil {
		return hcl2nix.Config{}, fmt.Errorf("Error fetching project name:", err)
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
	}, nil
}

func genPythonPoetryConf(pd *langdetect.ProjectDetails) hcl2nix.Config {
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
			Development: []string{"go@1.21.6", "gotools@0.16.1", "delve@1.22.0"},
			// todo: maybe we should dynamically inject the latest version of such runtime packages(cacert)?
			Runtime: []string{"cacert@3.95"},
		},
		GoModule: &hcl2nix.GoModule{
			Name:       name,
			SourcePath: entrypoint,
		},
	}

}
