package init

import (
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/langdetect"
)

func generatehcl2NixConf(pt langdetect.ProjectType, pd *langdetect.ProjectDetails) hcl2nix.Config {
	switch pt {
	case langdetect.GoModule:
		return genGoModuleConf(pd)
	case langdetect.PythonPoetry:
		return genPythonPoetryConf(pd)
	default:
		return hcl2nix.Config{
			Packages: hcl2nix.Packages{},
		}
	}
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
