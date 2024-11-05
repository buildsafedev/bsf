package init

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
)

var (
	commonDevDeps = []string{"coreutils-full@9.5", "bash@5.2.15"}
	commonRTDeps  = []string{"cacert@3.95"}
	rustDeps      = []string{"cargo@1.82.0"}
	pythonDeps    = []string{"python3@3.13.0", "poetry@1.8.4"}
	goDeps        = []string{"go@1.23.2", "gotools@0.25.0", "delve@1.23.1"}
	jsNpmDeps     = []string{"nodejs@23.1.0"}
)

func generatehcl2NixConf(pt langdetect.ProjectType, pd *langdetect.ProjectDetails, baseImgName string, addCommonDeps bool, commonDepsType string) (hcl2nix.Config, error) {
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
	case langdetect.BaseImage:
		return generateEmptyConf(baseImgName, addCommonDeps, commonDepsType), nil
	default:
		return hcl2nix.Config{
			Packages: hcl2nix.Packages{},
		}, fmt.Errorf("language is not supported")
	}
}

func generateEmptyConf(imageName string, addCommonDeps bool, commonDepsType string) hcl2nix.Config {
	devDeps := commonDevDeps
	if addCommonDeps {
		commonDepsType := strings.ToLower(strings.TrimSpace(commonDepsType))
		switch commonDepsType {
		case "go":
			devDeps = append(devDeps, goDeps...)
		case "python":
			devDeps = append(devDeps, pythonDeps...)
		case "rust":
			devDeps = append(devDeps, rustDeps...)
		case "jsnpm":
			devDeps = append(devDeps, jsNpmDeps...)
		}
	}
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: devDeps,
			Runtime:     commonRTDeps,
		},
		OCIArtifact: []hcl2nix.OCIArtifact{
			{
				Artifact: "pkgs",
				Name:     imageName,
				IsBase:   true,
				Layers: []string{
					"split(packages.runtime)",
					"split(packages.dev)",
				},
				Cmd:           []string{},
				Entrypoint:    []string{},
				EnvVars:       []string{},
				ExposedPorts:  []string{},
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

	rustDevDeps := append(commonDevDeps, rustDeps...)
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
	}, nil
}

func genPythonPoetryConf() hcl2nix.Config {
	// TODO: maybe we should note down the path of the poetry.lock file and use it here.
	poetryDevDeps := append(commonDevDeps, pythonDeps...)
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

	goDevDeps := append(commonDevDeps, goDeps...)
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

	nodeDevDeps := append(commonDevDeps, jsNpmDeps...)
	return hcl2nix.Config{
		Packages: hcl2nix.Packages{
			Development: nodeDevDeps,
			Runtime:     commonRTDeps,
		},
		JsNpmApp: &hcl2nix.JsNpmApp{
			PackageName: name,
			PackageRoot: "./.",
		},
	}, nil
}
